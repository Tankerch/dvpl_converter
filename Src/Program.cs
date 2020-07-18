using Force.Crc32;
using K4os.Compression.LZ4;
using System;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace DVPL_Converter
{
    class Program
    {
        // VARIABEL, PATH, AND CONST
        static readonly List<string> EXTENSION = new List<string> { ".dds", ".pvr", ".txt", ".fev", ".fsb", ".webp", ".yaml", ".xml", ".webp", ".mp3" };
        static readonly List<string> EXTENSION_00 = new List<string> { ".tex" };
        static readonly string PATH = Directory.GetCurrentDirectory(); //E:\Mods\TOOLS\GUP_SKINS_DVPL\MAIN ||  Directory.GetCurrentDirectory()


        static void Main()
        {
            PRINT_HEADER();

            bool DELETED_FILE_PERMISSION = GET_DELETE_CONFIRMATION();
            CONVERT_FILES(DELETED_FILE_PERMISSION);

            EXIT_PROGRAM();
        }

        static void PRINT_HEADER()
        {
            CONSOLE_COLOR(ConsoleColor.Cyan, "===================================================");
            CONSOLE_COLOR(ConsoleColor.Cyan, "DVPL_Converter C# Edition (https://github.com/Tankerch/DVPL_Converter)");
            CONSOLE_COLOR(ConsoleColor.Cyan, "Created By                    : Tankerch");
            CONSOLE_COLOR(ConsoleColor.Cyan, "Work Directory                : " + PATH);
            CONSOLE_COLOR(ConsoleColor.Cyan, "===================================================");
            CONSOLE_COLOR(ConsoleColor.Yellow, "Note : ONLY this extension will get converted :");
            Console.WriteLine("[{0} {1}]", string.Join(" ", EXTENSION.ToArray()), string.Join(" ", EXTENSION_00.ToArray()));

            Console.WriteLine("\nPress 'Enter' to continue");
            while (true)
            {
                if (Console.ReadKey().Key == ConsoleKey.Enter) break;
            }
        }

        static bool GET_DELETE_CONFIRMATION()
        {
            return false;
            Console.WriteLine("Do you want to delete the original files? Y/N");
            while (true)
            {
                string USER_INPUT = Console.ReadLine();
                if (USER_INPUT == "Y" || USER_INPUT == "y")
                {
                    return true;
                }
                else if (USER_INPUT == "N" || USER_INPUT == "n")
                {
                    return false;
                }
            }
        }

        static void CONVERT_FILES(bool DELETED_FILE_PERMISSION)
        {
            //BAD PRACTICE LOL. 
            string[] files = Directory.GetFiles(PATH, "*.*", SearchOption.AllDirectories);

            foreach (string file in files)
            {
                string FILE_EXTENSION = Path.GetExtension(file);
                if (EXTENSION.Contains(FILE_EXTENSION))
                {
                    CREATE_DVPL(LZ4Level.L03_HC, file);
                    HANDLE_CONVERTED(file, DELETED_FILE_PERMISSION);
                }
                else if (EXTENSION_00.Contains(FILE_EXTENSION))
                {
                    CREATE_DVPL(LZ4Level.L00_FAST, file);
                    HANDLE_CONVERTED(file, DELETED_FILE_PERMISSION);
                }
            }
        }

        static void CREATE_DVPL(LZ4Level COMPRESSION_TYPE, string file)
        {
            //Read original file to buffer
            byte[] ORIGINAL_DATA = File.ReadAllBytes(file);
            int ORIGINAL_SIZE = ORIGINAL_DATA.Length;


            // Calculate DVPL
            byte[] LZ4_CONTENT = new byte[LZ4Codec.MaximumOutputSize(ORIGINAL_SIZE)];
            int LZ4_SIZE = LZ4Codec.Encode(ORIGINAL_DATA, LZ4_CONTENT, COMPRESSION_TYPE);

            // SCREW THE IMPLEMENTATION OF LZ4 C# FOR LZ_00
            if (COMPRESSION_TYPE == LZ4Level.L00_FAST)
            {
                Buffer.BlockCopy(LZ4_CONTENT, 2, LZ4_CONTENT, 0, LZ4_CONTENT.Length - 2);
                LZ4_SIZE -= 2;
            }

            Array.Resize(ref LZ4_CONTENT, LZ4_SIZE);

            // Write DVPL byte to file
            byte[] DVPL_CONTENT = FORMAT_WG_DVPL(LZ4_CONTENT, LZ4_SIZE, ORIGINAL_SIZE, COMPRESSION_TYPE);
            File.WriteAllBytes(file + ".dvpl", DVPL_CONTENT);
        }

        static byte[] FORMAT_WG_DVPL(byte[] LZ4_CONTENT, int LZ4_SIZE, int ORIGINAL_SIZE, LZ4Level COMPRESSION_TYPE)
        {
            //Format from LZ4 compression to WG DVPL

            //Calculate CRC32, for what? IDK LOL ask WG
            uint LZ4_CRC32 = Crc32Algorithm.Compute(LZ4_CONTENT);

            byte[] DVPL_TEXT = Encoding.UTF8.GetBytes("DVPL");


            ushort COMPRESSION_TYPE_USHORT = (ushort)COMPRESSION_TYPE;

            if (COMPRESSION_TYPE != LZ4Level.L00_FAST) COMPRESSION_TYPE_USHORT -= 1;

            // BILL GATES FIX THIS MESS GODDAMMIT, I JUST WANNA TO APPEND THE BYTE FFS
            /*DVPL_FOOTER
             * {
             * (uint) ORIGINAL_SIZE,
             * (uint) LZ4_SIZE,
             * (uint) LZ4_CRC32,
             * (ushort) COMPRESSION_TYPE_SHORT,
             * (ushort) PADDING,
             * (DVPL_TEXT.Length) DVPL_TEXT}
            */

            byte[] DVPL_CONTENT = new byte[LZ4_CONTENT.Length + sizeof(uint) * 3 + sizeof(ushort) * 2 + DVPL_TEXT.Length];
            int OFFSET_ACCUMULATOR = 0;
            //LZ_CONTENT
            Buffer.BlockCopy(LZ4_CONTENT, 0, DVPL_CONTENT, OFFSET_ACCUMULATOR, LZ4_CONTENT.Length);
            OFFSET_ACCUMULATOR += LZ4_CONTENT.Length;
            //ORIGINAL_SIZE
            Buffer.BlockCopy(BitConverter.GetBytes((uint)ORIGINAL_SIZE), 0, DVPL_CONTENT, OFFSET_ACCUMULATOR, sizeof(uint));
            OFFSET_ACCUMULATOR += sizeof(uint);
            //LZ4_SIZE
            Buffer.BlockCopy(BitConverter.GetBytes((uint)LZ4_SIZE), 0, DVPL_CONTENT, OFFSET_ACCUMULATOR, sizeof(uint));
            OFFSET_ACCUMULATOR += sizeof(uint);
            //LZ4_CRC32
            Buffer.BlockCopy(BitConverter.GetBytes(LZ4_CRC32), 0, DVPL_CONTENT, OFFSET_ACCUMULATOR, sizeof(uint));
            OFFSET_ACCUMULATOR += sizeof(uint);
            //COMPRESSION_TYPE
            Buffer.BlockCopy(BitConverter.GetBytes(COMPRESSION_TYPE_USHORT), 0, DVPL_CONTENT, OFFSET_ACCUMULATOR, sizeof(ushort));
            OFFSET_ACCUMULATOR += sizeof(ushort);
            //PADDING
            Buffer.BlockCopy(BitConverter.GetBytes((ushort)0), 0, DVPL_CONTENT, OFFSET_ACCUMULATOR, sizeof(ushort));
            OFFSET_ACCUMULATOR += sizeof(ushort);
            //DVPL_TEXT
            Buffer.BlockCopy(DVPL_TEXT, 0, DVPL_CONTENT, OFFSET_ACCUMULATOR, DVPL_TEXT.Length);

            return DVPL_CONTENT;
        }


        static void HANDLE_CONVERTED(string FILE_PATH, bool DELETE_FILE_PERMISSION)
        {
            Console.WriteLine(FILE_PATH);
            if (DELETE_FILE_PERMISSION) File.Delete(FILE_PATH);
        }

        static void EXIT_PROGRAM()
        {
            Console.WriteLine("\nPress 'Enter' to exit program");
            while (true)
            {
                if (Console.ReadKey().Key == ConsoleKey.Enter) break;
            }
        }

        // Extension
        static void CONSOLE_COLOR(ConsoleColor color, string message)
        {
            Console.ForegroundColor = color;
            Console.WriteLine(message);
            Console.ResetColor();
        }
    }
}
