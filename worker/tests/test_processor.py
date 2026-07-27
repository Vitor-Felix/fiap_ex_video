import os
import stat
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from processor import process_video


class ProcessVideoTests(unittest.TestCase):
    def test_process_video_creates_zip_from_extracted_frames(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            base_dir = Path(tmpdir)
            upload_dir = base_dir / "uploads"
            output_dir = base_dir / "outputs"
            temp_dir = base_dir / "temp"
            upload_dir.mkdir(parents=True, exist_ok=True)
            output_dir.mkdir(parents=True, exist_ok=True)
            temp_dir.mkdir(parents=True, exist_ok=True)

            video_path = upload_dir / "sample.mp4"
            video_path.write_bytes(b"fake-video")

            fake_ffmpeg = base_dir / "bin" / "ffmpeg"
            fake_ffmpeg.parent.mkdir(parents=True, exist_ok=True)
            fake_ffmpeg.write_text(
                "#!/bin/sh\n"
                "set -e\n"
                "out_pattern=$6\n"
                'dir=$(dirname "$out_pattern")\n'
                'mkdir -p "$dir"\n'
                'touch "$dir/frame_0001.png"\n'
            )
            fake_ffmpeg.chmod(fake_ffmpeg.stat().st_mode | stat.S_IEXEC)

            os.environ["PATH"] = f"{fake_ffmpeg.parent}:{os.environ.get('PATH', '')}"
            os.environ["WORKER_BASE_DIR"] = str(base_dir)

            result = process_video(str(video_path), "video-123")

            self.assertTrue(result["success"])
            self.assertEqual(result["frame_count"], 1)
            self.assertTrue(result["zip_path"].endswith(".zip"))
            self.assertTrue((output_dir / result["zip_path"]).exists())


if __name__ == "__main__":
    unittest.main()
