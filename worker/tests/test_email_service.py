import os
import sys
import unittest
from unittest.mock import MagicMock, patch

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from email_service import send_error_email  # noqa: E402


class TestSendErrorEmail(unittest.TestCase):

    def setUp(self):
        """Garante que as variáveis SMTP estejam presentes em todos os testes."""
        os.environ["SMTP_HOST"] = "sandbox.smtp.mailtrap.io"
        os.environ["SMTP_PORT"] = "2525"
        os.environ["SMTP_USER"] = "test_user"
        os.environ["SMTP_PASSWORD"] = "test_pass"
        os.environ["SMTP_FROM"] = "worker@fiap-x.local"
        os.environ["SMTP_TO"] = "dev@fiap-x.local"

    @patch("email_service.smtplib.SMTP")
    def test_envia_email_quando_credenciais_presentes(self, mock_smtp_class):
        """Deve abrir conexão SMTP e chamar sendmail quando as credenciais existem."""
        mock_server = MagicMock()
        mock_smtp_class.return_value.__enter__ = MagicMock(return_value=mock_server)
        mock_smtp_class.return_value.__exit__ = MagicMock(return_value=False)

        send_error_email("video-001", "video.mp4", "FFmpeg exited with code 1")

        mock_smtp_class.assert_called_once_with("sandbox.smtp.mailtrap.io", 2525)
        mock_server.login.assert_called_once_with("test_user", "test_pass")
        mock_server.sendmail.assert_called_once()

        # Verifica que o destinatário correto foi usado
        args = mock_server.sendmail.call_args
        self.assertIn("dev@fiap-x.local", args[0][1])

    @patch("email_service.smtplib.SMTP")
    def test_nao_envia_quando_sem_credenciais(self, mock_smtp_class):
        """Não deve tentar conexão SMTP quando usuário ou senha estão vazios."""
        os.environ["SMTP_USER"] = ""
        os.environ["SMTP_PASSWORD"] = ""

        send_error_email("video-002", "video.mp4", "Erro qualquer")

        mock_smtp_class.assert_not_called()

    @patch("email_service.smtplib.SMTP")
    def test_nao_propaga_excecao_quando_smtp_falha(self, mock_smtp_class):
        """Falha no SMTP não deve estourar exceção para o chamador."""
        mock_server = MagicMock()
        mock_smtp_class.return_value.__enter__ = MagicMock(return_value=mock_server)
        mock_smtp_class.return_value.__exit__ = MagicMock(return_value=False)
        mock_server.sendmail.side_effect = Exception("connection refused")

        # Não deve lançar exceção
        try:
            send_error_email("video-003", "video.mp4", "Erro de processamento")
        except Exception as e:
            self.fail(f"send_error_email propagou exceção inesperadamente: {e}")

    @patch("email_service.smtplib.SMTP")
    def test_assunto_contem_nome_do_video(self, mock_smtp_class):
        """O assunto do e-mail deve incluir o nome do arquivo de vídeo."""
        import email as email_lib
        import email.header

        mock_server = MagicMock()
        mock_smtp_class.return_value.__enter__ = MagicMock(return_value=mock_server)
        mock_smtp_class.return_value.__exit__ = MagicMock(return_value=False)

        send_error_email("video-004", "meu_video_importante.mp4", "Crash")

        args = mock_server.sendmail.call_args
        raw_message = args[0][2]

        # Decodifica o cabeçalho Subject do MIME (pode estar em quoted-printable)
        parsed = email_lib.message_from_string(raw_message)
        subject_parts = email.header.decode_header(parsed["Subject"])
        subject = "".join(
            part.decode(enc or "utf-8") if isinstance(part, bytes) else part
            for part, enc in subject_parts
        )
        self.assertIn("meu_video_importante.mp4", subject)


if __name__ == "__main__":
    unittest.main()
