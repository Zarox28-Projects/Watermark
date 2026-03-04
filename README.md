<div align="center">
    <h1>Watermark</h1>
    <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white"/>
    <img src="https://img.shields.io/badge/License-AGPL%20v3-blue.svg?style=for-the-badge"/>
</div>

---

## Table of Contents

- [About](#about)
- [Installation](#installation)
  - [Homebrew](#homebrew)
  - [Build from Source](#build-from-source)
- [Usage](#usage)
  - [Command Line Options](#command-line-options)
- [License](#license)

## About

Watermark is an interactive command-line tool that simplifies the addition of watermarks to images and PDF files.

## Installation

### Homebrew

```bash
brew install zarox28/tap/watermark
```

### Build from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/Zarox28/Watermark.git
   cd Watermark
   ```

2. Build the application:

   ```bash
   go build -o ./build/watermark
   ```

3. (Optional) Add to your PATH for global access:
   ```bash
   sudo cp ./build/watermark /usr/local/bin/
   ```

## Usage

```bash
./watermark -i input.jpg -t "CONFIDENTIAL"
```

```bash
./watermark --input document.pdf --text "DRAFT - DO NOT DISTRIBUTE" --output secured_document.pdf
```

### Command Line Options

- `-i, --input` (required): Path to the input file (PNG, JPG, JPEG, or PDF)
- `-t, --text` (required): Text to use as watermark
- `-o, --output` (optional): Path for the output file (defaults to `output.ext`)

## License

This project is licensed under the AGPLv3 License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  Made with ❤️ by <a href="https://github.com/Zarox28">Zarox28</a>
</div>
