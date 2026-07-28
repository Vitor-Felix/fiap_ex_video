import os
import sys
import unittest
from unittest.mock import MagicMock, patch

# Garante que o diretório do worker esteja no path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

# psycopg2 não está disponível fora do container; mocka antes de importar db
sys.modules.setdefault("psycopg2", MagicMock())

from db import (  # noqa: E402
    update_status_to_completed,
    update_status_to_error,
    update_status_to_processing,
)


def make_mock_conn():
    """Retorna um mock de conexão psycopg2 pronto para uso."""
    mock_cursor = MagicMock()
    mock_conn = MagicMock()
    mock_conn.cursor.return_value.__enter__ = MagicMock(return_value=mock_cursor)
    mock_conn.cursor.return_value.__exit__ = MagicMock(return_value=False)
    return mock_conn, mock_cursor


class TestUpdateStatusToProcessing(unittest.TestCase):
    @patch("db.get_db_connection")
    def test_returns_true_on_success(self, mock_get_conn):
        mock_conn, mock_cursor = make_mock_conn()
        mock_get_conn.return_value = mock_conn

        result = update_status_to_processing("video-001")

        self.assertTrue(result)
        mock_conn.commit.assert_called_once()

    @patch("db.get_db_connection")
    def test_returns_false_on_db_error(self, mock_get_conn):
        mock_conn = MagicMock()
        mock_conn.cursor.side_effect = Exception("db offline")
        mock_get_conn.return_value = mock_conn

        result = update_status_to_processing("video-001")

        self.assertFalse(result)
        mock_conn.rollback.assert_called_once()


class TestUpdateStatusToCompleted(unittest.TestCase):
    @patch("db.get_db_connection")
    def test_returns_true_on_success(self, mock_get_conn):
        mock_conn, mock_cursor = make_mock_conn()
        mock_get_conn.return_value = mock_conn

        result = update_status_to_completed("video-002", "frames_video-002.zip", 42)

        self.assertTrue(result)
        mock_conn.commit.assert_called_once()

    @patch("db.get_db_connection")
    def test_returns_false_on_db_error(self, mock_get_conn):
        mock_conn = MagicMock()
        mock_conn.cursor.side_effect = Exception("connection reset")
        mock_get_conn.return_value = mock_conn

        result = update_status_to_completed("video-002", "frames.zip", 10)

        self.assertFalse(result)
        mock_conn.rollback.assert_called_once()


class TestUpdateStatusToError(unittest.TestCase):
    @patch("db.get_db_connection")
    def test_returns_true_on_success(self, mock_get_conn):
        mock_conn, mock_cursor = make_mock_conn()
        mock_get_conn.return_value = mock_conn

        result = update_status_to_error("video-003", "ffmpeg crash")

        self.assertTrue(result)
        mock_conn.commit.assert_called_once()

    @patch("db.get_db_connection")
    def test_returns_false_on_db_error(self, mock_get_conn):
        mock_conn = MagicMock()
        mock_conn.cursor.side_effect = Exception("timeout")
        mock_get_conn.return_value = mock_conn

        result = update_status_to_error("video-003", "ffmpeg crash")

        self.assertFalse(result)
        mock_conn.rollback.assert_called_once()


if __name__ == "__main__":
    unittest.main()
