#include <Windows.h>
#include <cstdio>

typedef struct _IDINFO {
    USHORT wGenConfig;
    USHORT wNumCyls;
    USHORT wReserved2;
    USHORT wNumHeads;
    USHORT wReserved4;
    USHORT wReserved5;
    USHORT wNumSectorsPerTrack;
    USHORT wVendorUnique[3];
    CHAR sSerialNumber[20];
    USHORT wBufferType;
    USHORT wBufferSize;
    USHORT wECCSize;
    CHAR sFirmwareRev[8];
    CHAR sModelNumber[40];
    USHORT wMoreVendorUnique;
    USHORT wReserved48;
    struct {
        USHORT reserved1 : 8;
        USHORT DMA : 1;
        USHORT LBA : 1;
        USHORT DisIORDY : 1;
        USHORT IORDY : 1;
        USHORT SoftReset : 1;
        USHORT Overlap : 1;
        USHORT Queue : 1;
        USHORT InlDMA : 1;
    } wCapabilities;
    USHORT wReserved1;
    USHORT wPIOTiming;
    USHORT wDMATiming;
    struct {
        USHORT CHSNumber : 1;
        USHORT CycleNumber : 1;
        USHORT UnltraDMA : 1;
        USHORT reserved : 13;
    } wFieldValidity;
    USHORT wNumCurCyls;
    USHORT wNumCurHeads;
    USHORT wNumCurSectorsPerTrack;
    USHORT wCurSectorsLow;
    USHORT wCurSectorsHigh;
    struct {
        USHORT CurNumber : 8;
        USHORT Multi : 1;
        USHORT reserved1 : 7;
    } wMultSectorStuff;
    ULONG dwTotalSectors;
    USHORT wSingleWordDMA;
    struct {
        USHORT Mode0 : 1;
        USHORT Mode1 : 1;
        USHORT Mode2 : 1;
        USHORT Reserved1 : 5;
        USHORT Mode0Sel : 1;
        USHORT Mode1Sel : 1;
        USHORT Mode2Sel : 1;
        USHORT Reserved2 : 5;
    } wMultiWordDMA;
    struct {
        USHORT AdvPOIModes : 8;
        USHORT reserved : 8;
    } wPIOCapacity;
    USHORT wMinMultiWordDMACycle;
    USHORT wRecMultiWordDMACycle;
    USHORT wMinPIONoFlowCycle;
    USHORT wMinPOIFlowCycle;
    USHORT wReserved69[11];
    struct {
        USHORT Reserved1 : 1;
        USHORT ATA1 : 1;
        USHORT ATA2 : 1;
        USHORT ATA3 : 1;
        USHORT ATA4 : 1;
        USHORT ATA5 : 1;
        USHORT ATA6 : 1;
        USHORT ATA7 : 1;
        USHORT ATA8 : 1;
        USHORT ATA9 : 1;
        USHORT ATA10 : 1;
        USHORT ATA11 : 1;
        USHORT ATA12 : 1;
        USHORT ATA13 : 1;
        USHORT ATA14 : 1;
        USHORT Reserved2 : 1;
    } wMajorVersion;
    USHORT wMinorVersion;
    USHORT wReserved82[6];
    struct {
        USHORT Mode0 : 1;
        USHORT Mode1 : 1;
        USHORT Mode2 : 1;
        USHORT Mode3 : 1;
        USHORT Mode4 : 1;
        USHORT Mode5 : 1;
        USHORT Mode6 : 1;
        USHORT Mode7 : 1;
        USHORT Mode0Sel : 1;
        USHORT Mode1Sel : 1;
        USHORT Mode2Sel : 1;
        USHORT Mode3Sel : 1;
        USHORT Mode4Sel : 1;
        USHORT Mode5Sel : 1;
        USHORT Mode6Sel : 1;
        USHORT Mode7Sel : 1;
    } wUltraDMA;
    USHORT wReserved89[167];
} IDINFO, *PIDINFO;


static void hexdump(const char *title, const void *pdata, int len) {
    printf("%s\n", title);
    int i, j, k, l;
    const char *data = (const char *) pdata;
    char buf[256], str[64], t[] = "0123456789ABCDEF";
    for (i = j = k = 0; i < len; i++) {
        if (0 == i % 16)
            j += sprintf(buf + j, "%08X  ", i);
        buf[j++] = t[0x0f & (data[i] >> 4)];
        buf[j++] = t[0x0f & data[i]];
        buf[j++] = ' ';
        str[k++] = isprint(data[i]) ? data[i] : '.';
        if (0 == (i + 1) % 16) {
            str[k] = 0;
            j += sprintf(buf + j, " |%s|\n", str);
            printf("%s", buf);
            j = k = buf[0] = str[0] = 0;
        }
    }
    str[k] = 0;
    if (k) {
        for (l = 0; l < 3 * (16 - k); l++)
            buf[j++] = ' ';
        j += sprintf(buf + j, " |%s|\n", str);
    }
    if (buf[0])
        printf("%s\n", buf);
    printf("\n");
}


void hexdump(void *addr, int size) {
    int skip = 0;
    int linechars = 16;
    int currline;
    int lasLine;
    int i;
    unsigned char *pc;
    char buff[256];
    char buff2[256];
    if (size - skip <= 0) {
        return;
    }
    lasLine = size / linechars;
    if (size % linechars != 0) {
        ++lasLine;
    }
    pc = (unsigned char *) addr;
    for (currline = skip / linechars; currline < lasLine; ++currline) {
        sprintf(buff, "  0x%04x ", currline * linechars);
        sprintf(buff2, "  ");
        for (i = 0; i < linechars; ++i) {
            int charno = currline * linechars + i;
            if (charno % 8 == 0) {
//                sprintf(buff, "  ");
            }
            if (charno >= skip && charno < size) {
                sprintf(buff, "%s %02x", buff, pc[charno]);
                if ((pc[charno] < 0x20) || (pc[charno] > 0x7e)) {
                    sprintf(buff2, "%s.", buff2);
                } else {
                    sprintf(buff2, "%s%c", buff2, pc[charno]);
                }
            } else {
                sprintf(buff, "%s __", buff);
                sprintf(buff2, "%s ", buff2);
            }
        }
        printf("%s%s\n", buff, buff2);
    }
}
int main_() {
    int i;
    unsigned char *buffer = static_cast<unsigned char *>(malloc(300));
    for (i = 0; i <= 255; ++i)
        buffer[i] = (unsigned char) i;
    hexdump(buffer, 300);
    free(buffer);
    return 0;
}

void exchange_char(char *in, char *out, size_t strlen_in) {
    for (size_t i = 0; i < (strlen_in); i += 2) {
        out[i] = in[i + 1];
        out[i + 1] = in[i];
    }
}

int main() {
    auto hDevice = CreateFileW(
            L"\\\\.\\PhysicalDrive0", GENERIC_READ | GENERIC_WRITE,
            FILE_SHARE_READ | FILE_SHARE_WRITE, NULL, OPEN_EXISTING, 0, NULL);
    if (hDevice == INVALID_HANDLE_VALUE) {
        // MessageBoxW(NULL, L"cccc", 0);
    }
    char OutBuffer[sizeof(SENDCMDOUTPARAMS) + IDENTIFY_BUFFER_SIZE - 1];
    GETVERSIONINPARAMS get_version;
    DWORD BytesReturned = 0;

    auto yy=sizeof(get_version);
    DeviceIoControl(hDevice, SMART_GET_VERSION, NULL, 0, &get_version,
                    sizeof(get_version), &BytesReturned, NULL);

    SENDCMDINPARAMS InBuffer = {0};
    InBuffer.irDriveRegs.bCommandReg =
            (get_version.bIDEDeviceMap & 0x10) ? ATAPI_ID_CMD : ID_CMD;


    auto lpInBuffer = InBuffer;
    auto InBufferSize = sizeof(SENDCMDINPARAMS) - 1;
      auto  b=reinterpret_cast<char*>(&InBuffer);
//    hexdump("lpInBuffer", &b[0], InBufferSize);
    hexdump( &b[0], InBufferSize);

    auto lpOutBuffer = OutBuffer;
    auto nOutBufferSize = sizeof(OutBuffer);

    DeviceIoControl(hDevice, SMART_RCV_DRIVE_DATA, &InBuffer,
                    sizeof(SENDCMDINPARAMS) - 1, OutBuffer,
                    sizeof(OutBuffer), &BytesReturned, NULL);
    auto out = reinterpret_cast<PSENDCMDOUTPARAMS>(OutBuffer);
    auto hd = reinterpret_cast<PIDINFO>(out->bBuffer);

    char data[512] = {};
    memcpy(data, out->bBuffer, 512);
    hexdump("data", data, 512);

    char disk_id[512] = {};
    char disk_model[512] = {};
    char sSerialNumber[512] = {};
    char sModelNumber[512] = {};
    exchange_char((hd->sSerialNumber), disk_id, sizeof(hd->sSerialNumber));
    exchange_char((hd->sModelNumber), disk_model, sizeof(hd->sModelNumber));

    exchange_char((hd->sSerialNumber), sSerialNumber, sizeof(hd->sSerialNumber));
    exchange_char((hd->sModelNumber), sModelNumber, sizeof(hd->sModelNumber));
    CloseHandle(hDevice);
    return 0;
}
