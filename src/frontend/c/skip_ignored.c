#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>

const uint8_t IGNORED_TOKEN_LEN = 7;
const uint8_t IGNORED_TOKEN[IGNORED_TOKEN_LEN] = {
    ' ', '\n', '\r', ',', '\t', '\v', '\f'
};
const uint8_t COMMENT_TOKEN = '#';
const uint8_t NEW_LINE_R    = '\r';
const uint8_t NEW_LINE_N    = '\n';

void skip_ignored(uint8_t *doc[], uint64_t docLen, uint64_t *ret[]){
    uint8_t c = 0;
    uint8_t hit = 0;
    for(uint64_t i = (*ret)[0]; i < docLen; i++){
        c = (*doc)[i];
        hit = 0;
        // skip ignored token
        for(int j = 0; j < IGNORED_TOKEN_LEN; j++){
            if(IGNORED_TOKEN[j] == c){
                hit = 1;
                break;
            }
        }
        // skip comment
        if(hit != 1 && c == COMMENT_TOKEN){
            for(;;){
                i ++;
                (*ret)[0] ++;
                if(i >= docLen){
                    return;
                }
                c = (*doc)[i];
                if(c == NEW_LINE_R || c == NEW_LINE_N){
                    hit = 1;
                    break;
                }
            }
        }
        // it is normal token, we need stop here
        if(hit != 1){
            return;
        }
        (*ret)[0] ++;
    }
    return;
}

