__constant uint K[64] = {
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
    0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
    0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
    0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
    0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
    0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2
};

inline uint rightRotate(uint value, uint count) {
    return (value >> count) | (value << (32 - count));
}

// SHA-256 funkcia
void sha256(const uchar* msg, int len, uint* hash) {
    uint h[8] = {
        0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
        0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19
    };

    uint w[64] = {0};
    for (int i = 0; i < len; i++) {
        w[i / 4] |= (uint)msg[i] << (24 - (i % 4) * 8);
    }
    w[len / 4] |= (uint)0x80 << (24 - (len % 4) * 8);
    w[15] = len * 8;

    for (int i = 16; i < 64; i++) {
        uint s0 = rightRotate(w[i - 15], 7) ^ rightRotate(w[i - 15], 18) ^ (w[i - 15] >> 3);
        uint s1 = rightRotate(w[i - 2], 17) ^ rightRotate(w[i - 2], 19) ^ (w[i - 2] >> 10);
        w[i] = w[i - 16] + s0 + w[i - 7] + s1;
    }

    uint a = h[0], b = h[1], c = h[2], d = h[3], e = h[4], f = h[5], g = h[6], h_val = h[7];

    for (int i = 0; i < 64; i++) {
        uint S1 = rightRotate(e, 6) ^ rightRotate(e, 11) ^ rightRotate(e, 25);
        uint ch = (e & f) ^ (~e & g);
        uint temp1 = h_val + S1 + ch + K[i] + w[i];
        uint S0 = rightRotate(a, 2) ^ rightRotate(a, 13) ^ rightRotate(a, 22);
        uint maj = (a & b) ^ (a & c) ^ (b & c);
        uint temp2 = S0 + maj;

        h_val = g;
        g = f;
        f = e;
        e = d + temp1;
        d = c;
        c = b;
        b = a;
        a = temp1 + temp2;
    }

    hash[0] = h[0] + a;
    hash[1] = h[1] + b;
    hash[2] = h[2] + c;
    hash[3] = h[3] + d;
    hash[4] = h[4] + e;
    hash[5] = h[5] + f;
    hash[6] = h[6] + g;
    hash[7] = h[7] + h_val;
}

// Funkcia na kontrolu počiatočných núl
bool hasLeadingZeros(uint* hash, uint numZero) {
    uint numBits = numZero * 4; // 1 hexadecimálna číslica = 4 bity

    for (uint i = 0; i < 8; i++) {
        for (int bit = 28; bit >= 0; bit -= 4) {
            uint nibble = (hash[i] >> bit) & 0xF;
            if (numBits == 0) return true;
            if (nibble != 0) return false;
            numBits -= 4;
        }
    }
    return true;
}
// Funkcia na konverziu čísla na reťazec
int itoa(int value, char* buffer, int base) {
    int i = 0;
    int isNegative = 0;

    // Ak je číslo záporné, zaznamenáme to a urobíme ho pozitívnym
    if (value < 0 && base == 10) {
        isNegative = 1;
        value = -value;
    }

    // Konvertujeme číslo do reťazca
    do {
        buffer[i++] = (char)(value % base + '0');
        value = value / base;
    } while (value);

    // Ak je číslo záporné, pridáme minus na začiatok
    if (isNegative) {
        buffer[i++] = '-';
    }

    buffer[i] = '\0';

    // Prevrátime reťazec, aby bol správne naformátovaný
    for (int j = 0; j < i / 2; j++) {
        char temp = buffer[j];
        buffer[j] = buffer[i - j - 1];
        buffer[i - j - 1] = temp;
    }

    return i;
}

// Funkcia na kombinovanie vstupného textu s číslom
int combineInputAndNumber(char* output, __global uchar* input, uint offset) {
    int len = 0;

    // Pridáme input
    while (input[len] != '\0') {
        output[len] = input[len];
        len++;
    }

    // Pridáme číslo (offset) ako reťazec
    char numberBuffer[16];  // Buffer na číslo, dostatočne veľký
    int numLen = itoa(offset, numberBuffer, 10);
    
    // Pridáme číslo k reťazcu
    for (int i = 0; i < numLen; i++) {
        output[len + i] = numberBuffer[i];
    }

    return len + numLen; // Celková dĺžka reťazca
}

__kernel void hashKernel(__global uchar* input, __global uint *offset, __global uint* output, __global uint *numZero, __global int* foundIndex) {
    int numStep = offset[1];
    int id = get_global_id(0);
    int numberStart = id + *offset;
    int numberGoal = numberStart + numStep;
    for(int i = numberStart; i<numberGoal; i++){
	if(*foundIndex!=-1){return;}
	uchar word[64] = {0};
	    int len = combineInputAndNumber((char*)word, input, i);

		    uint hash[8];
		    sha256(word, len, hash);

		    if (hasLeadingZeros(hash, *numZero)) {
		        // Skontrolujeme, či už bol nejaký výsledok uložený
		        if (atomic_cmpxchg(foundIndex, -1, id + *offset) == -1) {
		            // Prvý nájdený hash zapíšeme do output
		            output[0] = id + *offset; // Uložené číslo
		            for (int i = 0; i < 8; i++) {
		                output[i + 1] = hash[i]; // Uložený hash
		            }
		        }
		}
	}
}
