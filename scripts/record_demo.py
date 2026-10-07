#!/usr/bin/env python3
"""Run the real CLI; record stdout/timing and render a readable terminal replay.
Requires Pillow 11.3.0 only for rendering. No fabricated result strings.
"""
import hashlib
import json
import os
import subprocess
import textwrap
import time
from pathlib import Path
from PIL import Image, ImageDraw, ImageFont

root = Path(__file__).resolve().parents[1]
assets = root / "docs/assets"
binary = root / "bin/carebind"
start = time.monotonic()
process = subprocess.Popen([str(binary), "demo"], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
lines, events = [], []
for line in process.stdout:
    lines.append(line)
    events.append([round(time.monotonic() - start, 6), "o", line.replace("\n", "\r\n")])
assert process.wait() == 0, process.stderr.read()
transcript = "".join(lines)
(assets / "demo.txt").write_text(transcript)
header = {"version": 2, "width": 110, "height": 24, "title": "CareBind synthetic demo", "command": "bin/carebind demo"}
(assets / "demo.cast").write_text("\n".join(json.dumps(x) for x in [header, *events]) + "\n")
font_path = os.environ.get("CAREBIND_RECORD_FONT")
font = ImageFont.truetype(font_path, 18) if font_path else ImageFont.load_default(size=18)
small = ImageFont.load_default(size=16)
frames = []
blocks = transcript.strip().split("\n\n")[1:]
for index, block in enumerate(blocks):
    frame = Image.new("RGB", (1100, 400), "#0c1823")
    draw = ImageDraw.Draw(frame)
    draw.rounded_rectangle((20, 20, 1080, 380), radius=14, fill="#122636", outline="#294858", width=2)
    draw.text((44, 40), "carebind / source-bound evidence", font=font, fill="#87e3d3")
    draw.text((44, 82), "$ bin/carebind demo", font=font, fill="#f2f5f6")
    y = 125
    for line in block.splitlines():
        for wrapped in textwrap.wrap(line, width=93, replace_whitespace=False, drop_whitespace=False):
            draw.text((44, y), wrapped, font=font, fill="#edf3f7")
            y += 27
    draw.line((44, 298, 1055, 298), fill="#294858", width=1)
    draw.text((44, 314), "SYNTHETIC ONLY  |  physical identity: unverified  |  global freshness: unknown", font=small, fill="#ffc477")
    draw.text((44, 347), f"Actual CLI stdout · replay paced for reading · {index + 1}/{len(blocks)}", font=small, fill="#9ab5c6")
    frames.append(frame)
frames[1].save(assets / "demo.png")
frames[0].save(assets / "demo.gif", save_all=True, append_images=frames[1:], duration=3500, loop=0, optimize=False)
manifest = {"command": "bin/carebind demo", "exit_code": 0,
            "binary_sha256": hashlib.sha256(binary.read_bytes()).hexdigest(),
            "stdout_sha256": hashlib.sha256(transcript.encode()).hexdigest(),
            "capture": "actual process stdout, timestamped in demo.cast",
            "rendering": "Pillow 11.3.0 terminal replay, 3.5 seconds per case; not a desktop screenshot or physical trial",
            "source_sha256": {str(p.relative_to(root)): hashlib.sha256(p.read_bytes()).hexdigest()
                              for p in sorted(root.rglob("*.go")) if ".git" not in p.parts}}
(assets / "recording.json").write_text(json.dumps(manifest, indent=2) + "\n")
print("Recorded real CLI output and rendered 13-case replay")
