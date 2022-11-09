import json
from typing import List
from unittest.mock import ANY, Mock, call, patch

import pytest  # type: ignore

from raclambda.handler.raclambda_handler import (
    download_files,
    get_env_or_raise,
    get_new_files,
    handler,
    parse_event_message,
    upload_files,
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


def test_get_new_files():
    path_name = "path"
    file_names = ["foo", "bar"]
    with patch(
        "raclambda.handler.raclambda_handler.glob",
        return_value=["foo", "bar", "baz"],
    ):
        assert get_new_files(path_name, file_names) == ["baz"]


def test_upload_files():
    mocked_client = Mock()
    bucket_name = "bucket"
    path_name = "path"
    file_names = ["file1", "file2"]
    upload_files(mocked_client, bucket_name, path_name, file_names)
    mocked_client.upload_file.assert_has_calls([
        call('bucket', 'file1', 'path/file1'),
        call('bucket', 'file2', 'path/file2'),
    ])


def test_handle(monkeypatch):
    monkeypatch.setenv("RAC_PROJECT", "rac-project")
    monkeypatch.setenv("RAC_DREGS", "rac-dregs-bucket")
    rac_files = ["path/to/file.rac", "path/to/other-file.rac"]
    dregs_files = ["path/to/file.dregs", "path/to/other-file.dregs"]
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
    ) as patched_client, patch(
        'raclambda.handler.raclambda_handler.get_all_files',
        return_value=dregs_files,
    ) as patched_get_all, patch(
        'raclambda.handler.raclambda_handler.download_files',
    ) as patched_download, patch(
        'raclambda.handler.raclambda_handler.upload_files',
    ) as patched_upload, patch(
        'raclambda.handler.raclambda_handler.get_new_files',
        return_value=["path/to/new-file.dregs"]
    ) as patched_get_new, patch(
        'raclambda.handler.raclambda_handler.subprocess.call',
    ) as patched_call:
        patched_s3 = patched_client.return_value
        handler(event, None)
    patched_client.assert_called_once_with("s3")
    patched_download.assert_has_calls([
        call(patched_s3, "rac-bucket", ANY, rac_files),
        call(patched_s3, "rac-dregs-bucket", ANY, dregs_files),
    ])
    patched_get_all.assert_called_once_with("rac-dregs-bucket")
    patched_call.assert_called_once_with([
        "./rac",
        "-aws",
        "-project", "rac-project",
        "-dregs", ANY,
        ANY,
    ])
    patched_get_new.assert_called_once_with(ANY, dregs_files)
    patched_upload.assert_called_once_with(
        patched_s3,
        "rac-dregs-bucket",
        ANY,
        ["path/to/new-file.dregs"],
    )


def test_handle_raises(monkeypatch):
    monkeypatch.setenv("RAC_PROJECT", "rac-project")
    monkeypatch.setenv("RAC_DREGS", "rac-dregs-bucket")
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
    ) as patched_client:
        with pytest.raises(NothingToDo):
            handler(event, None)
    patched_client.assert_called_once_with("s3")
