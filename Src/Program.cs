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
        static readonly List<string> EXTENSION = new List<string>{ ".dds", ".pvr", ".txt", ".fev", ".fsb", ".webp", ".yaml", ".xml", ".webp" };
        static readonly List<string> EXTENSION_00 = new List<string> { ".tex"};
        static readonly string PATH = Directory.GetCurrentDirectory();


        static void Main()
        {
            PRINT_HEADER();

            Console.WriteLine("\nPress 'Enter' to continue");
            while (true)
            {
                if (Console.ReadKey().Key == ConsoleKey.Enter) break;
            }

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
            CONSOLE_COLOR(ConsoleColor.Yellow,"Note : ONLY this extension will get converted :" );
            Console.WriteLine("[{0} {1}]", string.Join(" ", EXTENSION.ToArray()), string.Join(" ", EXTENSION_00.ToArray()));
        }

        static bool GET_DELETE_CONFIRMATION()
        {
            return true;
        }

        static void CONVERT_FILES(bool DELETED_FILE_PERMISSION)
        {
            //BAD PRACTICE LOL. 
            string[] files =  Directory.GetFiles(PATH, "*.*", SearchOption.AllDirectories);

            foreach (string file in files)
            {

                string FILE_EXTENSION = Path.GetExtension(file);
                if (EXTENSION.Contains(FILE_EXTENSION))
                {
                    CREATE_DVPL(LZ4Level.L08_HC, file);
                    HANDLE_CONVERTED(file, DELETED_FILE_PERMISSION);
                }
                else if (EXTENSION_00.Contains(FILE_EXTENSION))
                {
                    CREATE_DVPL(LZ4Level.L12_MAX, file);
                    HANDLE_CONVERTED(file, DELETED_FILE_PERMISSION);
                }
            }
        }

        static void CREATE_DVPL(LZ4Level COMPRESSION_TYPE, string file)
        {
            //Read original file to buffer
            byte[] ORIGINAL_DATA = File.ReadAllBytes(file);
            int ORIGINAL_SIZE = ORIGINAL_DATA.Length;


            // Create DVPL file
            byte[] LZ4_CONTENT = new byte[LZ4Codec.MaximumOutputSize(ORIGINAL_SIZE)];
            int LZ4_SIZE = LZ4Codec.Encode(ORIGINAL_DATA, LZ4_CONTENT, COMPRESSION_TYPE);
            Array.Resize(ref LZ4_CONTENT, LZ4_SIZE);
            //File.WriteAllBytes(file + ".dvpl", FORMAT_WG_DVPL);
            byte[] DVPL_CONTENT = FORMAT_WG_DVPL(LZ4_CONTENT, LZ4_SIZE, ORIGINAL_SIZE, COMPRESSION_TYPE);
            File.WriteAllBytes(file + ".dvpl", DVPL_CONTENT);
        }

        static byte[] FORMAT_WG_DVPL(byte[] LZ4_CONTENT, int LZ4_SIZE, int ORIGINAL_SIZE, LZ4Level COMPRESSION_TYPE)
        {
            //Format from LZ4 compression to WG DVPL

            //Calculate CRC32, for what? IDK LOL ask WG
            //CRC32= crc32(LZ4_CONTENT) & 0xffffffff
            uint LZ4_CRC32 = Crc32Algorithm.Compute(LZ4_CONTENT);

            byte[] DVPL_TEXT = Encoding.UTF8.GetBytes("DVPL");
            byte DVPL_COMPRESION_TYPE = new byte();
            if (COMPRESSION_TYPE == LZ4Level.L08_HC)
            {
                DVPL_COMPRESION_TYPE = (byte)'0';
            }
            else
            {
                DVPL_COMPRESION_TYPE = (byte)'1'    ;
            }
            
            byte[] DVPL_CONTENT = new byte[sizeof(uint) + DVPL_TEXT.Length + sizeof(byte) + DVPL_TEXT.Length] ;
            Buffer.BlockCopy(LZ4_CONTENT, 0, DVPL_CONTENT, 0, LZ4_CONTENT.Length);
            Buffer.BlockCopy(BitConverter.GetBytes(ORIGINAL_SIZE) , 0, DVPL_CONTENT, LZ4_CONTENT.Length, sizeof(int));
            Buffer.BlockCopy(BitConverter.GetBytes(LZ4_SIZE), 0, DVPL_CONTENT, LZ4_CONTENT.Length, sizeof(int));
            return DVPL_CONTENT;
        }


        static void HANDLE_CONVERTED(string FILE_PATH, bool DELETE_FILE_PERMISSION)
        {
            Console.WriteLine(FILE_PATH);
            //if (DELETE_FILE_PERMISSION) ;
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
