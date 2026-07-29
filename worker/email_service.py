import os
import smtplib
from datetime import datetime
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText


def _get_smtp_config() -> dict:
    return {
        "host": os.getenv("SMTP_HOST", "sandbox.smtp.mailtrap.io"),
        "port": int(os.getenv("SMTP_PORT", "2525")),
        "user": os.getenv("SMTP_USER", ""),
        "password": os.getenv("SMTP_PASSWORD", ""),
        "from": os.getenv("SMTP_FROM", "worker@fiap-x.local"),
        "to": os.getenv("SMTP_TO", "dev@fiap-x.local"),
    }


def send_error_email(
    video_id: str,
    video_name: str,
    error_message: str,
    recipient: str = "",
) -> None:
    """
    Dispara um e-mail de alerta quando o processamento de um vídeo falha.

    O destinatário é resolvido na seguinte ordem:
      1. e-mail do usuário dono do vídeo (parâmetro recipient)
      2. fallback: variável de ambiente SMTP_TO
    Falhas no envio são apenas logadas — não interrompem o fluxo principal.
    """
    cfg = _get_smtp_config()

    if not cfg["user"] or not cfg["password"]:
        print("⚠️  [EMAIL] Credenciais SMTP não configuradas. Pulando envio.")
        return

    # Resolve destinatário: e-mail do usuário ou fallback de ambiente
    to_address = recipient.strip() if recipient and recipient.strip() else cfg["to"]
    if not to_address:
        print("⚠️  [EMAIL] Nenhum destinatário disponível. Pulando envio.")
        return

    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    subject = f"[FIAP-X] Falha no processamento do vídeo: {video_name}"

    body_plain = (
        f"O processamento do vídeo falhou.\n\n"
        f"ID:      {video_id}\n"
        f"Vídeo:   {video_name}\n"
        f"Erro:    {error_message}\n"
        f"Data:    {timestamp}\n"
    )

    body_html = f"""
    <html>
      <body style="font-family: sans-serif; color: #333;">
        <h2 style="color: #c0392b;">&#9888; Falha no processamento do vídeo</h2>
        <table cellpadding="8" cellspacing="0"
               style="border-collapse: collapse; width: 100%; max-width: 600px;">
          <tr style="background: #f9f9f9;">
            <td><strong>ID do Vídeo</strong></td>
            <td>{video_id}</td>
          </tr>
          <tr>
            <td><strong>Arquivo</strong></td>
            <td>{video_name}</td>
          </tr>
          <tr style="background: #f9f9f9;">
            <td><strong>Mensagem de Erro</strong></td>
            <td style="color: #c0392b;">{error_message}</td>
          </tr>
          <tr>
            <td><strong>Data / Hora</strong></td>
            <td>{timestamp}</td>
          </tr>
        </table>
        <p style="margin-top: 24px; font-size: 12px; color: #999;">
          FIAP-X Worker — notificação automática
        </p>
      </body>
    </html>
    """

    msg = MIMEMultipart("alternative")
    msg["Subject"] = subject
    msg["From"] = cfg["from"]
    msg["To"] = to_address
    msg.attach(MIMEText(body_plain, "plain"))
    msg.attach(MIMEText(body_html, "html"))

    try:
        with smtplib.SMTP(cfg["host"], cfg["port"]) as server:
            server.ehlo()
            server.starttls()
            server.login(cfg["user"], cfg["password"])
            server.sendmail(cfg["from"], [to_address], msg.as_string())
        print(f"📧 [EMAIL] Alerta enviado para {to_address} — vídeo {video_id}")
    except Exception as exc:
        # Nunca deixa o e-mail derrubar o fluxo principal
        print(f"❌ [EMAIL] Falha ao enviar alerta: {exc}")
