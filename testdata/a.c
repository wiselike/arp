#include <ctype.h>
#include <stdio.h>

#include <unistd.h>
#include <termios.h>
int ch;
struct termios old = {0};
void init() {
    if( tcgetattr(0, &old) < 0 ) perror("tcsetattr()");
    old.c_lflag &= ~ICANON;
    old.c_lflag &= ~ECHO;
    old.c_cc[VMIN] = 1;
    old.c_cc[VTIME] = 0;
    //if( tcsetattr(0, TCSANOW, &old) < 0 ) perror("tcsetattr ICANON");
}

void init2(){
    if(tcsetattr(0, TCSANOW, &old) < 0) perror("tcsetattr ~ICANON");
}

int main(void) {
    init();
    while (1) {
        init2();
        if( read(0, &ch,1) < 0 ) perror("read()");
        putchar((unsigned char)ch);
    }
    return 0;
}
