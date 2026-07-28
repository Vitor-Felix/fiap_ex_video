import os
import subprocess
import zipfile
from pathlib import Path
from typing import Dict


def get_base_dir() -> Path:
    return Path(os.getenv("WORKER_BASE_DIR", "/app")).resolve()


def process_video(video_path: str, video_id: str) -> Dict[str, object]:
    base_dir = get_base_dir()
    temp_dir = base_dir / "temp" / video_id
    output_dir = base_dir / "outputs"

    temp_dir.mkdir(parents=True, exist_ok=True)
    output_dir.mkdir(parents=True, exist_ok=True)

    frame_pattern = str(temp_dir / "frame_%04d.png")
    command = [
        "ffmpeg",
        "-i",
        video_path,
        "-vf",
        "fps=1",
        "-y",
        frame_pattern,
    ]

    try:
        result = subprocess.run(
            command,
            capture_output=True,
            text=True,
            check=True,
        )
        print(f"✅ [WORKER] FFmpeg concluído para vídeo {video_id}")
        print(result.stdout)
    except subprocess.CalledProcessError as exc:
        error_output = exc.stderr or exc.stdout or str(exc)
        raise RuntimeError(f"Erro ao executar ffmpeg: {error_output}") from exc

    frames = sorted(temp_dir.glob("*.png"))
    if not frames:
        raise RuntimeError("Nenhum frame foi extraído do vídeo")

    zip_filename = f"frames_{video_id}.zip"
    zip_path = output_dir / zip_filename

    try:
        with zipfile.ZipFile(zip_path, "w", compression=zipfile.ZIP_DEFLATED) as archive:
            for frame in frames:
                archive.write(frame, arcname=frame.name)
    except Exception as exc:
        raise RuntimeError(f"Erro ao criar arquivo ZIP: {exc}") from exc

    return {
        "success": True,
        "zip_path": zip_filename,
        "frame_count": len(frames),
    }
