#IN TESTING :
#- Sudah ditest untuk dx11.dds, dan XML
#- Sudah ditest untuk android

import lz4.frame
import zlib
import os
import msvcrt
from struct import *

extension = ['.dds', '.pvr', '.txt', '.tex', '.fev', '.fsb', '.webp', '.yaml', '.xml']
extension_00 = ['.webp', '.tex']

#konversi dari LZ4 ke DVPL
def buildDVPL(compressed, input_size, compressed_type = 2 ):
    compressed = compressed[19:-4] #Buang LZ4 Header and Footer

    compressed_size = len(compressed)
    CRC = zlib.crc32(compressed) & 0xffffffff #Hitung CRC-32 buat file terkompresi
    DVPL_text = bytearray('DVPL', 'utf-8')

    DVPL_footer = pack('LLLHH', input_size, compressed_size, CRC, compressed_type, 0) #input_size, compressed_size, new_CRC, compressed_type, padding
    return compressed + DVPL_footer + DVPL_text

#konversi dari LZ4 ke DVPL - garis besar
def writeDVPL(filename, compressed_type = 2):
    with open(filename, 'rb') as f: #Read file
        input_data = f.read()
        input_size = len(input_data)

    with open(filename + ".dvpl", 'wb') as f: #Tulis output
        if compressed_type == 0:
            compressed = bytearray(lz4.frame.compress(input_data, compression_level=8, block_size=lz4.frame.BLOCKSIZE_MAX4MB))
        else :
            compressed = bytearray(lz4.frame.compress(input_data, compression_level=12, block_size=lz4.frame.BLOCKSIZE_MAX4MB))
        f.write(buildDVPL(compressed,input_size,compressed_type))

#=====================================================================================================
#MAIN CODE HERE
path = os.getcwd()

print("WG PELER, BIKIN RIBET AJA DAH AMPE ADA DVPL SEGALA")
print("=========================================================")
print("Dibuat oleh : Tankerch (https://github.com/Tankerch/) ")
print("Area kerja direktori : " + path)
print("\nCatatan : HANYA file/subfolder dalam area kerja direktori dan extension ini yang akan DICONVERT :\n" + str(extension) )

print("Tekan tombol 'Enter' untuk melanjutkan")
while True:
    if msvcrt.getch() == b'\r':
        break

#Hapus file atau tidak?
jawaban = input("\nApakah anda ingin menghapus file sumber/asli? Y/N\n")
while True:
    if jawaban is 'Y' or jawaban is 'y' :
        DeleteConfirm = True
        break
    elif jawaban is 'N' or jawaban is 'n':
        DeleteConfirm = False
        break
    else:
        jawaban = input("\nApakah anda ingin menghapus file sumber/asli? Y/N\n")

#Path walking
StartCheckFile = True
files = []
# r=root, d=directories, f = files
for r, d, f in os.walk(path):
    for file in f:
        if file.endswith(tuple(extension)) :
            if StartCheckFile :
                print("\nList file yang dikonversi :")
                StartCheckFile = False
            paths = os.path.join(r, file)
            if file.endswith(tuple(extension_00)):
                writeDVPL(paths,0)
            else :
                writeDVPL(paths,2)
            print(paths)
            if DeleteConfirm :
                os.remove(paths)

if StartCheckFile :
    print("Tidak ada file yang dikonversi")

#End
print("Selesai, tekan 'Enter' untuk keluar dari aplikasi")
while True:
    if msvcrt.getch() == b'\r':
        break
