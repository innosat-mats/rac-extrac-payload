import json
from unittest.mock import call, patch

import pytest  # type: ignore
from raclambda.handler.raclambda_handler import (
    get_env_or_raise,
    handler,
    parse_event_message,
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


def test_handle(monkeypatch):
    monkeypatch.setenv("RAC_PROJECT", "rac-project")
    event = {
        "Records": [{
            "body": json.dumps({
                "objects": ["path/to/file.rac", "path/to/other-file.rac"],
                "bucket": "rac-bucket",
            })
        }]
    }

    with patch(
        'raclambda.handler.raclambda_handler.boto3.client',
    ) as patched_client, patch(
        'raclambda.handler.raclambda_handler.subprocess.call',
    ) as patched_call:
        patched_s3 = patched_client.return_value
        handler(event, None)
    patched_client.assert_called_once_with("s3")
    patched_s3.download_file.assert_has_calls([
        call(
            Bucket="rac-bucket",
            Key="path/to/file.rac",
            Filename="/tmp/file.rac",
        ),
        call(
            Bucket="rac-bucket",
            Key="path/to/other-file.rac",
            Filename="/tmp/other-file.rac",
        ),
    ])

    patched_call.assert_called_once_with([
        "./rac", "-aws", "-project", "rac-project", "/tmp/*.rac",
    ])
