/* 
 * File:   config.h
 * Author: Frank Breitinger
 *
 * Created on 1. Mai 2013, 12:15
 */

#ifndef MAIN_H
#define	MAIN_H

#include "../header/config.h"
#include "../header/hashing.h"
#include "../header/timing.h"
#include "../header/fingerprint.h"
#include "../header/fingerprintList.h"
#include "../header/bloomfilter.h"

// Global variable for the different modes
MODES *mode;
//FILE    *getFileHandle(char *filename);
void addPathToFingerprintList(FINGERPRINT_LIST *fpl, char *filename);
#endif	/* MAIN_H */

