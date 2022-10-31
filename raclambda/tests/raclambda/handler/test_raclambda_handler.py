import pytest  # type: ignore
from unittest.mock import Mock, patch

from raclambda.handler.raclambda_handler import (
    NoNewFiles,
    get_env_or_raise,
    handler,
)


def test_get_env_or_raise(monkeypatch):
    envvar = "ARANDOMENVKEY"
    monkeypatch.setenv(envvar, "envar")
    assert get_env_or_raise(envvar) == "envar"

    monkeypatch.delenv(envvar)
    with pytest.raises(ValueError, match=f"'{envvar}' not found in env"):
        get_env_or_raise(envvar) == "envar"


def test_handle(monkeypatch):
    monkeypatch.setenv("RAC_INPUT_BUCKET", "rac-bucket")
    monkeypatch.setenv("RAC_PROJECT", "rac-project")
    monkeypatch.setenv("RAC_QUEUE", "rac-queue")
    mocked_message = Mock()
    mocked_message.body = """{
        "Records": [{
            "s3": {
                "object": {
                    "key": "path/to/file.rac"
                }
            }
        }]
    }"""

    with patch(
        'raclambda.handler.raclambda_handler.boto3.resource',
    ) as patched_resource, patch(
        'raclambda.handler.raclambda_handler.boto3.client',
    ) as patched_client, patch(
        'raclambda.handler.raclambda_handler.get_messages',
        return_value=[mocked_message],
    ) as patched_get_messages, patch(
        'raclambda.handler.raclambda_handler.subprocess.call',
    ) as patched_call:
        patched_queue = patched_resource.return_value
        patched_queue.get_queue_by_name.return_value = "queue"
        patched_s3 = patched_client.return_value
        handler(None, None)
    patched_resource.assert_called_once_with("sqs")
    patched_queue.get_queue_by_name.assert_called_once_with(
        QueueName="rac-queue",
    )
    patched_get_messages.assert_called_once_with("queue")
    patched_client.assert_called_once_with("s3")
    patched_s3.download_file.assert_called_once_with(
        Bucket="rac-bucket",
        Key="path/to/file.rac",
        Filename="/tmp/file.rac",
    )
    patched_call.assert_called_once_with([
        "./rac", "-aws", "-project", "rac-project", "/tmp/*.rac",
    ])
    mocked_message.delete.assert_called_once()


def test_handler_raises_no_files(monkeypatch):
    monkeypatch.setenv("RAC_INPUT_BUCKET", "rac-bucket")
    monkeypatch.setenv("RAC_PROJECT", "rac-project")
    monkeypatch.setenv("RAC_QUEUE", "rac-queue")

    with patch(
        'raclambda.handler.raclambda_handler.boto3.resource',
    ) as patched_resource, patch(
        'raclambda.handler.raclambda_handler.get_messages',
        return_value=[],
    ):
        patched_queue = patched_resource.return_value
        patched_queue.get_queue_by_name.return_value = "queue"
        with pytest.raises(
            NoNewFiles,
            match="Got no new files from rac-queue",
        ):
            handler(None, None)
