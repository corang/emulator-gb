package cart

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type RomHeader struct {
	EntryPoint   [4]byte
	NintendoLogo [48]byte

	Title           [16]byte
	LicenseCodeNew  [2]byte
	SGBFlag         byte
	CartridgeType   byte
	ROMSize         byte
	RAMSize         byte
	DestinationCode byte
	LicenseCodeOld  byte
	Version         byte
	HeaderChecksum  byte
	GlobalChecksum  [2]byte
}

type Cart struct {
	Size   uint32
	Data   []byte
	Header *RomHeader
}

var RomTypes = map[byte]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x04: "0x04 ???",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x07: "0x07 ???",
	0x08: "ROM+RAM 1",
	0x09: "ROM+RAM+BATTERY 1",
	0x0A: "0x0A ???",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0E: "0x0E ???",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY 2",
	0x11: "MBC3",
	0x12: "MBC3+RAM 2",
	0x13: "MBC3+RAM+BATTERY 2",
	0x14: "0x14 ???",
	0x15: "0x15 ???",
	0x16: "0x16 ???",
	0x17: "0x17 ???",
	0x18: "0x18 ???",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x1F: "0x1F ???",
	0x20: "MBC6",
	0x21: "0x21 ???",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
}

var NewLicenseCodes = map[string]string{
	"00": "None",
	"01": "Nintendo R&D1",
	"08": "Capcom",
	"13": "Electronic Arts",
	"18": "Hudson Soft",
	"19": "b-ai",
	"20": "kss",
	"22": "pow",
	"24": "PCM Complete",
	"25": "san-x",
	"28": "Kemco Japan",
	"29": "seta",
	"30": "Viacom",
	"31": "Nintendo",
	"32": "Bandai",
	"33": "Ocean/Acclaim",
	"34": "Konami",
	"35": "Hector",
	"37": "Taito",
	"38": "Hudson",
	"39": "Banpresto",
	"41": "Ubi Soft",
	"42": "Atlus",
	"44": "Malibu",
	"46": "angel",
	"47": "Bullet-Proof",
	"49": "irem",
	"50": "Absolute",
	"51": "Acclaim",
	"52": "Activision",
	"53": "American sammy",
	"54": "Konami",
	"55": "Hi tech entertainment",
	"56": "LJN",
	"57": "Matchbox",
	"58": "Mattel",
	"59": "Milton Bradley",
	"60": "Titus",
	"61": "Virgin",
	"64": "LucasArts",
	"67": "Ocean",
	"69": "Electronic Arts",
	"70": "Infogrames",
	"71": "Interplay",
	"72": "Broderbund",
	"73": "sculptured",
	"75": "sci",
	"78": "THQ",
	"79": "Accolade",
	"80": "misawa",
	"83": "lozc",
	"86": "Tokuma Shoten Intermedia",
	"87": "Tsukuda Original",
	"91": "Chunsoft",
	"92": "Video system",
	"93": "Ocean/Acclaim",
	"95": "Varie",
	"96": "Yonezawa/s’pal",
	"97": "Kaneko",
	"99": "Pack in soft",
	"A4": "Konami (Yu-Gi-Oh!)",
}

var OldLicenseCodes = map[byte]string{
	0x00: "None",
	0x01: "Nintendo",
	0x08: "Capcom",
	0x09: "Hot-B",
	0x0A: "Jaleco",
	0x0B: "Coconuts Japan",
	0x0C: "Elite Systems",
	0x13: "EA (Electronic Arts)",
	0x18: "Hudsonsoft",
	0x19: "ITC Entertainment",
	0x1A: "Yanoman",
	0x1D: "Japan Clary",
	0x1F: "Virgin Interactive",
	0x24: "PCM Complete",
	0x25: "San-X",
	0x28: "Kotobuki Systems",
	0x29: "Seta",
	0x30: "Infogrames",
	0x31: "Nintendo",
	0x32: "Bandai",
	0x33: "Indicates that the New licensee code should be used instead.",
	0x34: "Konami",
	0x35: "HectorSoft",
	0x38: "Capcom",
	0x39: "Banpresto",
	0x3C: ".Entertainment i",
	0x3E: "Gremlin",
	0x41: "Ubisoft",
	0x42: "Atlus",
	0x44: "Malibu",
	0x46: "Angel",
	0x47: "Spectrum Holoby",
	0x49: "Irem",
	0x4A: "Virgin Interactive",
	0x4D: "Malibu",
	0x4F: "U.S. Gold",
	0x50: "Absolute",
	0x51: "Acclaim",
	0x52: "Activision",
	0x53: "American Sammy",
	0x54: "GameTek",
	0x55: "Park Place",
	0x56: "LJN",
	0x57: "Matchbox",
	0x59: "Milton Bradley",
	0x5A: "Mindscape",
	0x5B: "Romstar",
	0x5C: "Naxat Soft",
	0x5D: "Tradewest",
	0x60: "Titus",
	0x61: "Virgin Interactive",
	0x67: "Ocean Interactive",
	0x69: "EA (Electronic Arts)",
	0x6E: "Elite Systems",
	0x6F: "Electro Brain",
	0x70: "Infogrames",
	0x71: "Interplay",
	0x72: "Broderbund",
	0x73: "Sculptered Soft",
	0x75: "The Sales Curve",
	0x78: "t.hq",
	0x79: "Accolade",
	0x7A: "Triffix Entertainment",
	0x7C: "Microprose",
	0x7F: "Kemco",
	0x80: "Misawa Entertainment",
	0x83: "Lozc",
	0x86: "Tokuma Shoten Intermedia",
	0x8B: "Bullet-Proof Software",
	0x8C: "Vic Tokai",
	0x8E: "Ape",
	0x8F: "I’Max",
	0x91: "Chunsoft Co.",
	0x92: "Video System",
	0x93: "Tsubaraya Productions Co.",
	0x95: "Varie Corporation",
	0x96: "Yonezawa/S’Pal",
	0x97: "Kaneko",
	0x99: "Arc",
	0x9A: "Nihon Bussan",
	0x9B: "Tecmo",
	0x9C: "Imagineer",
	0x9D: "Banpresto",
	0x9F: "Nova",
	0xA1: "Hori Electric",
	0xA2: "Bandai",
	0xA4: "Konami",
	0xA6: "Kawada",
	0xA7: "Takara",
	0xA9: "Technos Japan",
	0xAA: "Broderbund",
	0xAC: "Toei Animation",
	0xAD: "Toho",
	0xAF: "Namco",
	0xB0: "acclaim",
	0xB1: "ASCII or Nexsoft",
	0xB2: "Bandai",
	0xB4: "Square Enix",
	0xB6: "HAL Laboratory",
	0xB7: "SNK",
	0xB9: "Pony Canyon",
	0xBA: "Culture Brain",
	0xBB: "Sunsoft",
	0xBD: "Sony Imagesoft",
	0xBF: "Sammy",
	0xC0: "Taito",
	0xC2: "Kemco",
	0xC3: "Squaresoft",
	0xC4: "Tokuma Shoten Intermedia",
	0xC5: "Data East",
	0xC6: "Tonkinhouse",
	0xC8: "Koei",
	0xC9: "UFL",
	0xCA: "Ultra",
	0xCB: "Vap",
	0xCC: "Use Corporation",
	0xCD: "Meldac",
	0xCE: ".Pony Canyon or",
	0xCF: "Angel",
	0xD0: "Taito",
	0xD1: "Sofel",
	0xD2: "Quest",
	0xD3: "Sigma Enterprises",
	0xD4: "ASK Kodansha Co.",
	0xD6: "Naxat Soft",
	0xD7: "Copya System",
	0xD9: "Banpresto",
	0xDA: "Tomy",
	0xDB: "LJN",
	0xDD: "NCS",
	0xDE: "Human",
	0xDF: "Altron",
	0xE0: "Jaleco",
	0xE1: "Towa Chiki",
	0xE2: "Yutaka",
	0xE3: "Varie",
	0xE5: "Epcoh",
	0xE7: "Athena",
	0xE8: "Asmik ACE Entertainment",
	0xE9: "Natsume",
	0xEA: "King Records",
	0xEB: "Atlus",
	0xEC: "Epic/Sony Records",
	0xEE: "IGS",
	0xF0: "A Wave",
	0xF3: "Extreme Entertainment",
	0xFF: "LJN",
}

func getLicense(oldCode byte, newCode [2]byte) string {
	var license string

	if oldCode == 0x33 {
		license = NewLicenseCodes[string(newCode[:])]
	} else {
		license = OldLicenseCodes[oldCode]
	}

	if license == "" {
		license = "Unknown"
	}
	return license
}

func Load(filename string) (*Cart, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	headerData := fileContents[0x100:]

	var header RomHeader

	reader := bytes.NewReader(headerData)

	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		fmt.Println("Error decoding header:", err)
		return nil, err
	}

	fmt.Println("Title:", string(header.Title[:]))
	fmt.Println("Cartridge Type:", RomTypes[header.CartridgeType])
	fmt.Println("ROM Size:", header.ROMSize)
	fmt.Println("RAM Size:", header.RAMSize)
	fmt.Println("License:", getLicense(header.LicenseCodeOld, header.LicenseCodeNew))
	fmt.Println("Version:", header.Version)

	var checksum byte = 0

	for chkAddr := uint16(0x0134); chkAddr < 0x014D; chkAddr++ {
		checksum = checksum - fileContents[chkAddr] - 1
	}

	if checksum == header.HeaderChecksum {
		fmt.Println("Checksum: PASSED")
	} else {
		fmt.Println("Checksum: FAILED")
	}
	
	loadedCart := &Cart{
		Size:   uint32(len(fileContents)),
		Data:   fileContents,
		Header: &header,
	}

	return loadedCart, nil
}
