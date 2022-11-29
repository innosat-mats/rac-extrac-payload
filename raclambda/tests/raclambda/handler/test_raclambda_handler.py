import json
from pathlib import Path
from typing import List
from unittest.mock import ANY, Mock, call, patch

import botocore
import pytest  # type: ignore
from botocore.stub import Stubber

from raclambda.handler.raclambda_handler import (
    download_files,
    format_rclone_command,
    get_env_or_raise,
    get_rclone_config_path,
    handler,
    parse_event_message,
    NothingToDo,
)


def test_get_env_or_raise(monkeypatch):
    envvar = "ARANDOMENVKEY"
    monkeypatch.setenv(envvar, "envar")
    assert get_env_or_raise(envvar) == "envar"

    monkeypatch.delenv(envvar)
    with pytest.raises(ValueError, match=f"'{envvar}' not found in env"):
        get_env_or_raise(envvar) == "envar"


def test_parse_event_message():
    event = {
        "Records": [
            {
                "body": '{"bucket": "rac-bucket", "objects": ["path/to/file.rac"]}'  # noqa: E501
            }
        ]
    }
    assert parse_event_message(event) == (["path/to/file.rac"], "rac-bucket")


def test_download_files(tmp_path):
    mocked_client = Mock()
    bucket_name = "bucket"
    file_names = ["file1", "file2"]
    download_files(mocked_client, bucket_name, tmp_path, file_names)
    mocked_client.download_file.assert_has_calls([
        call('bucket', 'file1', f'{tmp_path}/file1'),
        call('bucket', 'file2', f'{tmp_path}/file2'),
    ])


def test_handler(monkeypatch):
    monkeypatch.setenv("RAC_PROJECT", "rac-project")
    monkeypatch.setenv("RAC_DREGS", "rac-dregs-bucket")
    monkeypatch.setenv("RAC_OUTPUT", "rac-output-bucket")
    monkeypatch.setenv("RCLONE_CONFIG_SSM_NAME", "rclone-config")
    rac_files = ["path/to/file.rac", "path/to/other-file.rac"]
    event = {
        "Records": [{
            "body": json.dumps({
                "objects": rac_files,
                "bucket": "rac-bucket",
            })
        }]
    }

    with patch(
        'raclambda.handler.raclambda_handler.boto3.client',
    ) as patched_boto, patch(
        'raclambda.handler.raclambda_handler.get_rclone_config_path',
        return_value="/rclone/config",
    ) as patched_rclone_config, patch(
        'raclambda.handler.raclambda_handler.download_files',
    ) as patched_download, patch(
        'raclambda.handler.raclambda_handler.subprocess.call',
    ) as patched_call:
        patched_client = patched_boto.return_value
        handler(event, None)

    patched_boto.assert_has_calls([call("s3"), call("ssm")], any_order=False)
    patched_download.assert_called_once_with(
        patched_client, "rac-bucket", ANY, rac_files,
    )
    patched_rclone_config.assert_called_once_with(
        patched_client, "rclone-config",
    )
    patched_call.assert_has_calls([
        call(["rclone", "--config", "/rclone/config", "copy", "S3:rac-dregs-bucket", ANY, "--size-only"]),  # noqa: E501
        call(["./rac", "-parquet", "-project", ANY, "-dregs", ANY, ANY]),
        call(["rclone", "--config", "/rclone/config", "copy", ANY, "S3:rac-output-bucket"]),  # noqa: E501
        call(["rclone", "--config", "/rclone/config", "copy", ANY, "S3:rac-dregs-bucket", "--size-only"]),  # noqa: E501
    ], any_order=False)


def test_handler_raises_nothing_to_do(monkeypatch):
    monkeypatch.setenv("RAC_PROJECT", "rac-project")
    monkeypatch.setenv("RAC_DREGS", "rac-dregs-bucket")
    monkeypatch.setenv("RAC_OUTPUT", "rac-output-bucket")
    monkeypatch.setenv("RCLONE_CONFIG_SSM_NAME", "rclone-config")
    rac_files: List[str] = []
    event = {
        "Records": [{
            "body": json.dumps({
                "objects": rac_files,
                "bucket": "rac-bucket",
            })
        }]
    }

    with patch(
        'raclambda.handler.raclambda_handler.boto3.client',
    ) as patched_boto:
        with pytest.raises(NothingToDo):
            handler(event, None)

    patched_boto.assert_called_once_with("s3")


def test_rclone_config_path():
    ssm_parameter = "param"

    ssm_client = botocore.session.get_session().create_client(
        "ssm",
        region_name="eu-north-1"
    )
    stubber = Stubber(ssm_client)
    stubber.add_response(
        "get_parameter",
        {"Parameter": {"Value": "config"}},
        expected_params={"Name": ssm_parameter, "WithDecryption": True}
    )
    stubber.activate()

    name = get_rclone_config_path(ssm_client, ssm_parameter)

    path = Path(name)
    assert path.exists()
    assert path.read_text() == "config"
    path.unlink()


def test_format_rclone_command_sloppy():
    assert format_rclone_command("config", "from_path", "to_path", True) == [
        "rclone", "--config", "config", "copy", "from_path", "to_path", "--size-only",  # noqa: E501
    ]


def test_format_rclone_command():
    assert format_rclone_command("config", "from_path", "to_path") == [
        "rclone", "--config", "config", "copy", "from_path", "to_path",
    ]
