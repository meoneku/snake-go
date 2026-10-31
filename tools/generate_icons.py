import os
import sys
from pathlib import Path
from PIL import Image
from playwright.sync_api import sync_playwright

SVG_CONTENT = '''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <defs>
    <linearGradient id="bgGrad" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#1b2338"/>
      <stop offset="100%" stop-color="#0b0f1a"/>
    </linearGradient>
    <linearGradient id="snakeGrad" x1="0%" y1="100%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#14b8a6"/>
      <stop offset="100%" stop-color="#4ade80"/>
    </linearGradient>
    <linearGradient id="appleGrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#ef4444"/>
      <stop offset="100%" stop-color="#991b1b"/>
    </linearGradient>
    <filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="8" result="blur" />
      <feComposite in="SourceGraphic" in2="blur" operator="over" />
    </filter>
  </defs>

  <!-- Squircle Base Canvas (App Icon Shape) -->
  <rect width="512" height="512" rx="112" fill="url(#bgGrad)"/>

  <!-- Subtle Cyber Grid Dots -->
  <g fill="#ffffff" fill-opacity="0.12">
    <circle cx="96" cy="96" r="5"/>
    <circle cx="256" cy="96" r="5"/>
    <circle cx="416" cy="96" r="5"/>
    <circle cx="96" cy="256" r="5"/>
    <circle cx="416" cy="256" r="5"/>
    <circle cx="96" cy="416" r="5"/>
    <circle cx="256" cy="416" r="5"/>
    <circle cx="416" cy="416" r="5"/>
  </g>

  <!-- Cyan Neon Border Ring -->
  <rect width="504" height="504" x="4" y="4" rx="108" fill="none" stroke="#38bdf8" stroke-width="6" stroke-opacity="0.45"/>

  <!-- Snake Body (Curved loop from bottom left to center) -->
  <path d="M 76 360 C 90 440, 220 450, 250 370 C 270 300, 190 260, 180 190" 
        fill="none" stroke="#052e1b" stroke-width="76" stroke-linecap="round"/>
  <path d="M 76 360 C 90 440, 220 450, 250 370 C 270 300, 190 260, 180 190" 
        fill="none" stroke="url(#snakeGrad)" stroke-width="64" stroke-linecap="round"/>

  <!-- Snake Head -->
  <circle cx="190" cy="176" r="64" fill="#052e1b"/>
  <circle cx="190" cy="176" r="56" fill="#4ade80"/>

  <!-- Tongue (Forked, flickering out towards the apple) -->
  <path d="M 240 192 L 288 208 L 304 196 M 288 208 L 304 220" 
        fill="none" stroke="#f472b6" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"/>

  <!-- Eyes & Pupils (Looking at the apple with high-shine highlights) -->
  <!-- Right Eye (Closer) -->
  <circle cx="212" cy="148" r="17" fill="#ffffff"/>
  <circle cx="218" cy="148" r="9" fill="#052e1b"/>
  <circle cx="215" cy="144" r="3.5" fill="#ffffff"/>

  <!-- Left Eye -->
  <circle cx="168" cy="144" r="15" fill="#ffffff"/>
  <circle cx="172" cy="144" r="8" fill="#052e1b"/>
  <circle cx="170" cy="141" r="3" fill="#ffffff"/>

  <!-- Juicy Red Apple -->
  <g transform="translate(324, 224)">
    <!-- Stem -->
    <rect x="-4" y="-72" width="10" height="24" rx="4" fill="#924d0e"/>
    <!-- Leaf -->
    <ellipse cx="20" cy="-66" rx="20" ry="12" fill="#4ade80" transform="rotate(-20 20 -66)"/>
    <!-- Shadow Base -->
    <circle cx="0" cy="6" r="64" fill="#991b1b"/>
    <!-- Apple Body -->
    <circle cx="0" cy="0" r="60" fill="url(#appleGrad)"/>
    <!-- Gloss Highlight -->
    <ellipse cx="-18" cy="-20" rx="18" ry="12" fill="#ffffff" fill-opacity="0.45" transform="rotate(-30 -18 -20)"/>
  </g>

  <!-- Golden Sparkle Stars (+50 bonus hint) -->
  <!-- Big Star -->
  <g transform="translate(408, 96)">
    <path d="M 0 -28 Q 0 0 28 0 Q 0 0 0 28 Q 0 0 -28 0 Q 0 0 0 -28 Z" fill="#fbbf24"/>
    <circle cx="0" cy="0" r="4" fill="#ffffff"/>
  </g>
  <!-- Small Star -->
  <g transform="translate(112, 108)">
    <path d="M 0 -14 Q 0 0 14 0 Q 0 0 0 14 Q 0 0 -14 0 Q 0 0 0 -14 Z" fill="#fbbf24" fill-opacity="0.85"/>
  </g>
</svg>
'''

def main():
    root_dir = Path(__file__).resolve().parent.parent
    assets_dir = root_dir / "gui" / "assets"
    assets_dir.mkdir(parents=True, exist_ok=True)

    # 1. Save master SVG
    svg_path = assets_dir / "icon.svg"
    svg_path.write_text(SVG_CONTENT.strip(), encoding="utf-8")
    print(f"Master SVG saved to {svg_path}")

    # 2. Render 1024x1024 Master PNG via Playwright (Headless Chromium)
    master_png = assets_dir / "icon_master_1024.png"
    print("Rendering high-res master PNG via Chromium...")
    
    html_content = f'''<!DOCTYPE html>
<html>
<head>
  <style>
    * {{ margin: 0; padding: 0; }}
    body {{ background: transparent; overflow: hidden; }}
    svg {{ width: 1024px; height: 1024px; }}
  </style>
</head>
<body>
  {SVG_CONTENT}
</body>
</html>'''

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page(viewport={"width": 1024, "height": 1024}, device_scale_factor=1)
        page.set_content(html_content)
        # Take transparent PNG screenshot
        page.screenshot(path=str(master_png), omit_background=True)
        browser.close()

    print(f"Master PNG rendered: {master_png}")

    # 3. Resample into all standard sizes using Pillow LANCZOS
    base_img = Image.open(master_png).convert("RGBA")
    
    sizes = [512, 256, 128, 64, 48, 32, 16]
    generated_images = {}
    for sz in sizes:
        resized = base_img.resize((sz, sz), Image.Resampling.LANCZOS)
        out_png = assets_dir / f"icon_{sz}.png"
        resized.save(out_png, "PNG", optimize=True)
        generated_images[sz] = resized
        print(f"Generated: {out_png} ({sz}x{sz})")

    # Standard icon.png (256x256)
    default_icon = assets_dir / "icon.png"
    generated_images[256].save(default_icon, "PNG", optimize=True)
    print(f"Generated default: {default_icon}")

    # 4. Generate multi-resolution Windows .ico
    ico_path = root_dir / "gui" / "icon.ico"
    # Windows .ico includes [16, 32, 48, 64, 128, 256]
    ico_sizes = [(sz, sz) for sz in [16, 32, 48, 64, 128, 256]]
    # Pillow save with sizes takes the base image and resizes or uses pre-resized
    base_img.save(
        ico_path,
        format="ICO",
        sizes=ico_sizes
    )
    print(f"Generated Windows multi-resolution ICO: {ico_path}")

    # Clean up temporary master if desired, or keep as high-res reference
    if master_png.exists():
        master_png.unlink()

    print("All icon assets generated successfully!")

if __name__ == "__main__":
    main()
