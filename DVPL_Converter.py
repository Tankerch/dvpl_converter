import enum
import os
from msvcrt import getch
from struct import pack
from sys import exit
from zlib import crc32

import lz4.frame

# Extension, compression type & vars
EXTENSION: list = ['.dds', '.pvr', '.txt', '.fev',
                   '.fsb', '.webp', '.yaml', '.xml', '.webp']
EXTENSION_00: list = ['.tex']

PATH: str = os.getcwd()


class CompressionType(enum.Enum):
    HIGH_COMPRESSION: int = 1
    MAX_COMPRESSION: int = 2

# MAIN CODE


def MAIN():
    PRINT_HEADER()

    print("Press 'Enter' to continue")
    while True:
        if getch() == b'\r':
            break

    DELETED_FILE_PERMISSION: bool = GET_DELETE_CONFIRMATION()
    CONVERT_FILES(DELETED_FILE_PERMISSION)

    EXIT_PROGRAM()


def PRINT_HEADER() -> None:
    print("=========================================================")
    print("DVPL_Converter Python Edition    https://github.com/Tankerch/DVPL_Converter")
    print("Created By       : Tankerch (https://github.com/Tankerch/) ")
    print("Work Directory   : " + PATH)
    print("=========================================================")
    print("Note : ONLY this extension will get converted :\n" +
          str(EXTENSION + EXTENSION_00))


def GET_DELETE_CONFIRMATION() -> bool:
    while True:
        keystroke: str = input(
            "\nDo you want to delete the original files? Y/N\n")
        if keystroke == 'y' or keystroke == 'Y':
            return True
        elif keystroke == 'n' or keystroke == 'N':
            return False


def CONVERT_FILES(DELETE_FILE_PERMISSION: bool) -> None:
    # Path Walking
    for root, directories, files in os.walk(PATH):
        for file in files:
            FILE_PATH: str = os.path.join(root, file)
            if file.endswith(tuple(EXTENSION)):
                CREATE_DVPL(CompressionType.HIGH_COMPRESSION, FILE_PATH)
                HANDLE_CONVERTED(FILE_PATH, DELETE_FILE_PERMISSION)
            elif file.endswith(tuple(EXTENSION_00)):
                CREATE_DVPL(CompressionType.MAX_COMPRESSION, FILE_PATH)
                HANDLE_CONVERTED(FILE_PATH, DELETE_FILE_PERMISSION)


def CREATE_DVPL(COMPRESION_TYPE: int, FILE_PATH: str) -> None:
    '''Create DVPL from original file
        - Compression type == HIGH COMPRESSION -> COMPRESSION LEVEL 8
        - Compression type == MAX COMPRESSION -> COMPRESSION LEVEL 12'''
    # Read original file to buffer
    with open(FILE_PATH, 'rb') as file:
        ORIGINAL_DATA: bytes = file.read()
        ORIGINAL_SIZE: int = len(ORIGINAL_DATA)

    # Create DVPL file
    with open(FILE_PATH + ".dvpl", 'wb') as file:
        if COMPRESION_TYPE == CompressionType.HIGH_COMPRESSION:
            LZ4_CONTENT: bytearray = bytearray(lz4.frame.compress(
                ORIGINAL_DATA, compression_level=8, block_size=lz4.frame.BLOCKSIZE_MAX4MB))
        else:
            LZ4_CONTENT: bytearray = bytearray(lz4.frame.compress(
                ORIGINAL_DATA, compression_level=12, block_size=lz4.frame.BLOCKSIZE_MAX4MB))
        file.write(FORMAT_WG_DVPL(LZ4_CONTENT, ORIGINAL_SIZE, COMPRESION_TYPE))


def FORMAT_WG_DVPL(LZ4_CONTENT: bytearray, ORIGINAL_SIZE: int, COMPRESION_TYPE: int) -> bytearray:
    ''' Format from LZ4 compression to WG DVPL '''
    LZ4_CONTENT: bytearray = LZ4_CONTENT[19:-4]  # Strip LZ4 Header and Footer
    LZ4_SIZE: int = len(LZ4_CONTENT)

    # Calculate CRC32, for what? IDK LOL ask WG
    CRC32: int = crc32(LZ4_CONTENT) & 0xffffffff

    DVPL_TEXT: bytearray = bytearray('DVPL', 'utf-8')
    DVPL_FOOTER: bytes = pack('LLLHH', ORIGINAL_SIZE, LZ4_SIZE,
                              CRC32, COMPRESION_TYPE, 0)

    return LZ4_CONTENT + DVPL_FOOTER + DVPL_TEXT


def HANDLE_CONVERTED(FILE_PATH: str, DELETE_FILE_PERMISSION: bool) -> None:
    ''' Print converted files to console. if DELETE_FILE_PERMISSION == true, then delete the original file  '''
    print(FILE_PATH)
    if DELETE_FILE_PERMISSION:
        os.remove(FILE_PATH)


def EXIT_PROGRAM() -> None:
    print("Press 'Enter' to exit")
    while True:
        if getch() == b'\r':
            exit(0)


if __name__ == "__main__":
    MAIN()
