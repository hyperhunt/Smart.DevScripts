
// GO Lang :: SmartGo/Tests :: Smart.Go.Framework
// (c) 2020 unix-world.org
// r.20200507.1905 :: STABLE

package main

import (
	"os"
	"log"
	"fmt"
	smart "github.com/unix-world/smartgo"
	uid   "github.com/unix-world/smartgo/uuid"
)

const (
	THE_TPL = `[%%%COMMENT%%%]
This is comment one ...
[%%%/COMMENT%%%]
Hallo, this is Markers TPL:[%%%|N%%%][###MARKER|json###][%%%|T%%%][###MARKER2|url|html###][%%%|T%%%][###MARKER2|js###]
[%%%COMMENT%%%]
This is another comment ...
[%%%/COMMENT%%%]`
)


func LogToConsoleWithColors() {
	//--
//	smart.LogToFile("DEBUG", "test-mtpl.log", false) // plain log
//	smart.LogToFile("DEBUG", "test-mtpl.json.log", true) // json log
	smart.LogToConsole("DEBUG", true)
	//--
	log.Println("[DEBUG] Debugging")
	log.Println("[NOTICE] Notice")
	log.Println("[WARNING] Warning")
	log.Println("[ERROR] Error")
	log.Println("A log message, with no type")
	//--
}

func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION



func testDTimeNow() {

	dTimeStr := ""

	dtObjUtc := smart.DateTimeStructUtc(dTimeStr)
	fmt.Println("Date Obj Json", smart.JsonEncode(dtObjUtc))
	if(dtObjUtc.Status != "OK") {
		fatalError("testDTimeNow ERROR (UTC): " + dtObjUtc.ErrMsg)
	} //end if

	dtObjLoc := smart.DateTimeStructLocal(dTimeStr)
	fmt.Println("Date Obj Json", smart.JsonEncode(dtObjLoc))
	if(dtObjLoc.Status != "OK") {
		fatalError("testDTimeNow ERROR (LOCAL): " + dtObjLoc.ErrMsg)
	} //end if

} //END FUNCTION


func testDTimeParse() {

	dTimeStr := "2020-05-07 08:03:07 +0300"

	fmt.Println("Input Date:", dTimeStr)

	dtObjUtc := smart.DateTimeStructUtc(dTimeStr)
	fmt.Println("Date Obj Json", smart.JsonEncode(dtObjUtc))
	if(dtObjUtc.Status != "OK") {
		fatalError("testDTimeParse ERROR (UTC): " + dtObjUtc.ErrMsg)
	} //end if

	dtObjLoc := smart.DateTimeStructLocal(dTimeStr)
	fmt.Println("Date Obj Json", smart.JsonEncode(dtObjLoc))
	if(dtObjLoc.Status != "OK") {
		fatalError("testDTimeParse ERROR (LOCAL): " + dtObjLoc.ErrMsg)
	} //end if

} //END FUNCTION


func testAllTrimWhitespaces(input string) {

	input = smart.StrTrimWhitespaces(input)

	if(input == "") {
		fatalError("TRIM Test Input is Empty !")
	} //end if

	var whiteSpaces string = "\r\n\n\t\r\t \x0B \x00"

	var trimTest string = ""

	trimTest = smart.StrTrimWhitespaces(whiteSpaces + input + whiteSpaces)
	fmt.Println("TRIM Test: `" + trimTest + "`")
	if(trimTest != input) {
		fatalError("TRIM Test ERROR !")
	} //end if

	trimTest = smart.StrTrimLeftWhitespaces(whiteSpaces + input)
	fmt.Println("TRIM LEFT Test: `" + trimTest + "`")
	if(trimTest != input) {
		fatalError("LEFT TRIM Test ERROR !")
	} //end if

	trimTest = smart.StrTrimRightWhitespaces(input + whiteSpaces)
	fmt.Println("TRIM RIGHT Test: `" + trimTest + "`")
	if(trimTest != input) {
		fatalError("RIGHT TRIM Test ERROR !")
	} //end if

} //END FUNCTION


func testStrings() {

	var test1Str string = "1234567890 abcdefgh șȘțȚâÂăĂîÎ _ abcdefghijklmnopqrstuvwxyz-ABCDEFGHIJKLMNOPQRSTUVWXYZ:;\"'~`!@#$%^&*()+=[]{}|\\<>,.?/\t\r\n@  abcdefgh șȘțȚâÂăĂîÎ üÜöÖäÄ"

	if(smart.StrMBSubstr(test1Str, 0, 0) != smart.StrSubstr(test1Str, 0, 0)) {
		fatalError("SubString Comparison Test FAILED !")
	} //end if

	var needleStr string = ""
	var testThePos int = -777

	needleStr = "abcdefgh șȘțȚâÂăĂîÎ"
	testThePos = smart.StrPos(test1Str, needleStr)
	if(testThePos != 11) {
		fatalError("StrPos (case sensitive) Test FAILED. Expected result is 11 but get: ", testThePos)
	} //end if

	needleStr = "abcdefgh șȘțȚÂÂăĂîÎ"
	testThePos = smart.StrPos(test1Str, needleStr)
	if(testThePos != -1) {
		fatalError("StrPos (case sensitive) Test FAILED. Expected result is -1 but get: ", testThePos)
	} //end if

	needleStr = "abcdefgh ȘșțȚâÂăĂîÎ"
	testThePos = smart.StrIPos(test1Str, needleStr)
	if(testThePos != 11) {
		fatalError("StrIPos (case insensitive) Test FAILED. Expected result is 11 but get: ", testThePos)
	} //end if

	needleStr = "abcdefgh  ȘșțȚâÂăĂîÎ"
	testThePos = smart.StrIPos(test1Str, needleStr)
	if(testThePos != -1) {
		fatalError("StrIPos (case insensitive) Test FAILED. Expected result is -1 but get: ", testThePos)
	} //end if

	needleStr = "abcdefgh șȘțȚâÂăĂîÎ"
	testThePos = smart.StrRPos(test1Str, needleStr)
	if(testThePos != 122) {
		fatalError("StrRPos (case sensitive) Test FAILED. Expected result is 122 but get: ", testThePos)
	} //end if

	needleStr = "abcdefgh șȘțȚÂÂăĂîÎ"
	testThePos = smart.StrRPos(test1Str, needleStr)
	if(testThePos != -1) {
		fatalError("StrRPos (case sensitive) Test FAILED. Expected result is -1 but get: ", testThePos)
	} //end if

	needleStr = "abcdefgh ȘșțȚâÂăĂîÎ"
	testThePos = smart.StrRIPos(test1Str, needleStr)
	if(testThePos != 122) {
		fatalError("StrRIPos (case insensitive) Test FAILED. Expected result is 122 but get: ", testThePos)
	} //end if

	needleStr = "abcdefgh  ȘșțȚâÂăĂîÎ"
	testThePos = smart.StrRIPos(test1Str, needleStr)
	if(testThePos != -1) {
		fatalError("StrRIPos (case insensitive) Test FAILED. Expected result is -1 but get: ", testThePos)
	} //end if

} //END FUNCTION


func testFileSystem(input string) {

	var thePath string = "a/path/to/some/file.ext"

	fmt.Println("------------------------- Path Get BaseName / DirName TESTS -------------------------")

	fBaseName := smart.PathBaseName(thePath)
	if(fBaseName != "file.ext") {
		fatalError("Path Get BaseName FAILED ; BaseName = `" + fBaseName + "` ; Path = `" + thePath + "`")
	} //end if
	fDirName := smart.PathDirName(thePath)
	if(fDirName != "a/path/to/some") {
		fatalError("Path Get DirName FAILED ; DirName = `" + fDirName + "` ; Path = `" + thePath + "`")
	} //end if
	fmt.Println("Path Get BaseName and DirName OK ; Path = `" + thePath + "` ; BaseName = `" + fBaseName + "` ; DirName = `" + fDirName + "`")

	fmt.Println("------------------------- Absolute Path TESTS -------------------------")

	if(smart.PathIsAbsolute(thePath) == true) {
		fatalError("Absolute Path Detection Failed. This is not an absolute path: " + thePath)
	} //end if
	fmt.Println("Absolute Path Test OK (not absolute path): " + thePath)

	var theAPath string = "/" + thePath
	if(smart.PathIsAbsolute(theAPath) != true) {
		fatalError("Absolute Path Detection Failed. This is an absolute path: " + theAPath)
	} //end if
	fmt.Println("Absolute Path Test OK (it is an absolute path): " + theAPath)

	theAPath = ":" + thePath
	if(smart.PathIsAbsolute(theAPath) != true) {
		fatalError("Absolute Path Detection Failed. This is an absolute path: " + theAPath)
	} //end if
	fmt.Println("Absolute Path Test OK (it is an absolute path): " + theAPath)

	theAPath = "C:" + thePath
	if(smart.PathIsAbsolute(theAPath) != true) {
		fatalError("Absolute Path Detection Failed. This is an absolute path: " + theAPath)
	} //end if
	fmt.Println("Absolute Path Test OK (it is an absolute path): " + theAPath)

	fmt.Println("------------------------- Backward Unsafe Path TESTS -------------------------")

	if(smart.PathIsBackwardUnsafe(thePath) == true) {
		fatalError("Backward Unsafe Path Detection Failed. This is not a backward unsafe path: " + thePath)
	} //end if
	fmt.Println("Backward Unsafe Path Test OK (not backward unsafe path): " + thePath)

	theAPath = thePath + "/../"
	if(smart.PathIsBackwardUnsafe(theAPath) != true) {
		fatalError("Backward Unsafe Path Detection Failed. This is a backward unsafe path: " + theAPath)
	} //end if
	fmt.Println("Backward Unsafe Path Test OK (it is a backward unsafe path): " + theAPath)

	theAPath = thePath + "/./"
	if(smart.PathIsBackwardUnsafe(theAPath) != true) {
		fatalError("Backward Unsafe Path Detection Failed. This is a backward unsafe path: " + theAPath)
	} //end if
	fmt.Println("Backward Unsafe Path Test OK (it is a backward unsafe path): " + theAPath)

	theAPath = thePath + "/.."
	if(smart.PathIsBackwardUnsafe(theAPath) != true) {
		fatalError("Backward Unsafe Path Detection Failed. This is a backward unsafe path: " + theAPath)
	} //end if
	fmt.Println("Backward Unsafe Path Test OK (it is a backward unsafe path): " + theAPath)

	theAPath = thePath + "../"
	if(smart.PathIsBackwardUnsafe(theAPath) != true) {
		fatalError("Backward Unsafe Path Detection Failed. This is a backward unsafe path: " + theAPath)
	} //end if
	fmt.Println("Backward Unsafe Path Test OK (it is a backward unsafe path): " + theAPath)

	fmt.Println("------------------------- Path/File Exists TESTS -------------------------")

	if(smart.PathExists("./nonexisting") == true) {
		fatalError("Errors encountered while testing a non-existing path exists")
	} //end if
	if(smart.PathIsDir("./nonexisting") == true) {
		fatalError("Errors encountered while testing a non-existing path is a dir")
	} //end if
	if(smart.PathIsFile("./nonexisting") == true) {
		fatalError("Errors encountered while testing a non-existing path is a file")
	} //end if

	if(smart.PathExists("./") != true) {
		fatalError("Errors encountered while testing if the current directory ./ exists")
	} //end if
	if(smart.PathIsDir("./") != true) {
		fatalError("Errors encountered while testing if the current directory ./ is a dir")
	} //end if
	if(smart.PathIsFile("./") == true) {
		fatalError("Errors encountered while testing if the current directory ./ is not a file")
	} //end if

	if(smart.PathExists("mtpl.go") != true) {
		fatalError("Errors encountered while testing if the file mtpl.go exists")
	} //end if
	if(smart.PathIsDir("mtpl.go") == true) {
		fatalError("Errors encountered while testing if the file mtpl.go is not a dir")
	} //end if
	if(smart.PathIsFile("mtpl.go") != true) {
		fatalError("Errors encountered while testing if the file mtpl.go is a file")
	} //end if

	fmt.Println("Path Exists Test Result: PASSED");

	fmt.Println("------------------------- Dir Create and Rename TESTS -------------------------")

	var theBaseDir string = "test-Dir/"
	var theDir string = theBaseDir

	if(smart.PathExists(theDir)) {
		fatalError("Dir Tests EXISTS. Remove it manually before starting the tests: `" + theDir)
	} //end if

	isOkDirCreate1, errmsg := smart.SafePathDirCreate(theDir, false, false)
	if(errmsg != "") {
		fatalError("Errors encountered while creating the test dir `" + theDir + "`:", errmsg)
	} //end if
	if(isOkDirCreate1 != true) {
		fatalError("Failed to create the test dir (level 1) `" + theDir + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; DirCreate: `" + theDir + "` ; Result =", isOkDirCreate1)

	var theRenamedDir string = theDir + "LEVEL2/level3-X/"
	theDir = theDir + "LEVEL2/level3/"
	isOkDirCreate2, errmsg := smart.SafePathDirCreate(theDir, false, false)
	if((errmsg == "") || (isOkDirCreate2 == true)) {
		fatalError("Errors encountered while creating the recursive test dir (levels ++) without allowed (should return error but was not) `" + theDir + "`: ErrMsg=", errmsg, " ; isSuccess=", isOkDirCreate2)
	} //end if
	fmt.Println("Test PASSED ; DirCreateRecursive/Dissalowed: `" + theDir + "` ; Result =", isOkDirCreate2)

	isOkDirCreate3, errmsg := smart.SafePathDirCreate(theDir, true, false)
	if(errmsg != "") {
		fatalError("Errors encountered while creating the test recursive dir (levels ++) `" + theDir + "`:", errmsg)
	} //end if
	if(isOkDirCreate3 != true) {
		fatalError("Failed to create the test recursive dir (levels ++) `" + theDir + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; DirCreateRecursive: `" + theDir + "` ; Result =", isOkDirCreate3)

	isOkRenameDir, errmsg := smart.SafePathDirRename(theDir, theRenamedDir, false)
	if(errmsg != "") {
		fatalError("Errors encountered while renaming the test dir `" + theDir + "` to `" + theRenamedDir + "`:", errmsg)
	} //end if
	if(isOkRenameDir != true) {
		fatalError("Failed to rename the test dir `" + theDir + "` to `" + theRenamedDir + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; DirRename: `" + theDir + "` to `" + theRenamedDir + "` ; Result =", isOkRenameDir)
	theDir = theRenamedDir

	if((!smart.PathExists(theDir)) || (!smart.PathIsDir(theDir))) {
		fatalError("Dir Tests FAILED. Cannot find the Path/Dir: `" + theDir)
	} //end if
	fmt.Println("Tests PASSED ; DirCreate/DirRename [ +/- Recursive ]: `" + theDir + "` ; ALL")

	fmt.Println("------------------------- File Write TESTS -------------------------")

	var theTestFile string = theBaseDir + "test-FILE.txt"
	var theFContents string = input + "\n" + smart.DateNowLocal() + "\n"

	var theEmptyTestFile string = theTestFile + ".empty.TXT"
	isESuccess, errmsg := smart.SafePathFileWrite("", "w", theEmptyTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while writing the empty test file `" + theEmptyTestFile + "`:", errmsg)
	} //end if
	if(isESuccess != true) {
		fatalError("Failed to write the empty test file `" + theEmptyTestFile + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; WriteEmptyFile: `" + theEmptyTestFile + "` ; Result =", isESuccess)

	isSuccess, errmsg := smart.SafePathFileWrite(theFContents, "w", theTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while writing the test file `" + theTestFile + "`:", errmsg)
	} //end if
	if(isSuccess != true) {
		fatalError("Failed to write the test file `" + theTestFile + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; WriteFile: `" + theTestFile + "` ; Result =", isSuccess)

	fmt.Println("------------------------- File Write/Append TESTS -------------------------")

	isASuccess, errmsg := smart.SafePathFileWrite(theFContents, "a", theTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while write/append the test file `" + theTestFile + "`:", errmsg)
	} //end if
	if(isASuccess != true) {
		fatalError("Failed to write/append the test file `" + theTestFile + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; Append/WriteFile: `" + theTestFile + "` ; Result =", isASuccess)

	fmt.Println("------------------------- File Read TESTS -------------------------")

	fEmptyContent, errmsg := smart.SafePathFileRead(theEmptyTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while reading the empty test file `" + theEmptyTestFile + "`:", errmsg)
	} //end if
	if(fEmptyContent != "") {
		fatalError("Failed to read the empty test file `" + theEmptyTestFile + "` (not empty content)")
	} //end if
	fmt.Println("Test PASSED ; ReadEmptyFile `" + theEmptyTestFile + "` ; FileLength =", len(fEmptyContent), "bytes")

	fContent, errmsg := smart.SafePathFileRead(theTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while reading the test file `" + theTestFile + "`:", errmsg)
	} //end if
	if(fContent == "") {
		fatalError("Failed to read the test file `" + theTestFile + "` (empty content)")
	} //end if
	if(fContent != (theFContents + theFContents)) {
		fatalError("Failed to read the test file `" + theTestFile + "` (content is wrong): `" + fContent + "` instead of `" + theFContents + "`")
	} //end if
	fmt.Println("Test PASSED ; ReadFile `" + theTestFile + "` ; FileLength =", len(fContent), "bytes")

	fmt.Println("------------------------- File Rename TESTS -------------------------")

	var theNewTestFile string = theDir + "test-Renamed.txt"

	isReSuccess, errmsg := smart.SafePathFileRename(theTestFile, theNewTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while renaming the test file `" + theTestFile + "` to `" + theNewTestFile + "`:", errmsg)
	} //end if
	if(isReSuccess != true) {
		fatalError("Failed to rename the test file `" + theTestFile + "` to `" + theNewTestFile + "` (result is FALSE)")
	} //end if
	if((smart.PathExists(theTestFile) != false) || (smart.PathIsFile(theTestFile) != false) || (smart.PathIsDir(theTestFile) != false)) {
		fatalError("Failed to rename the test file `" + theTestFile + "` to `" + theNewTestFile + "` (old path still exists after rename ...)")
	} //end if
	if((smart.PathExists(theNewTestFile) != true) || (smart.PathIsFile(theNewTestFile) != true)) {
		fatalError("Failed to rename the test file `" + theTestFile + "` to `" + theNewTestFile + "` (new file does not exists after rename ...)")
	} //end if
	fmt.Println("Test PASSED ; RenameFile: `" + theTestFile + "` to `" + theNewTestFile + "` ; Result =", isReSuccess)

	fmt.Println("------------------------- File Copy TESTS -------------------------")

	theCopyEmptyTestFile := theEmptyTestFile + "_copied-from-ORIGINAL-EMPTY_File-here.txt"
	isECpSuccess, errmsg := smart.SafePathFileCopy(theEmptyTestFile, theCopyEmptyTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while copying the empty test file `" + theEmptyTestFile + "` to `" + theCopyEmptyTestFile + "`:", errmsg)
	} //end if
	if(isECpSuccess != true) {
		fatalError("Failed to copy the empty test file `" + theEmptyTestFile + "` to `" + theCopyEmptyTestFile + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; CopyEmptyFile: `" + theEmptyTestFile + "` to `" + theCopyEmptyTestFile + "` ; Result =", isECpSuccess)

	theCopyTestFile := theBaseDir + "test-renamed-from-ORIGINAL-and-after-Copied-here.txt"
	isCpSuccess, errmsg := smart.SafePathFileCopy(theNewTestFile, theCopyTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while copying the test file `" + theNewTestFile + "` to `" + theCopyTestFile + "`:", errmsg)
	} //end if
	if(isCpSuccess != true) {
		fatalError("Failed to copy the test file `" + theNewTestFile + "` to `" + theCopyTestFile + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; CopyFile: `" + theNewTestFile + "` to `" + theCopyTestFile + "` ; Result =", isCpSuccess)

	fmt.Println("------------------------- File Delete TESTS -------------------------")

	isDSuccess, errmsg := smart.SafePathFileDelete(theNewTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while deleting the test file `" + theNewTestFile + "`:", errmsg)
	} //end if
	if(isDSuccess != true) {
		fatalError("Failed to delete the test file `" + theNewTestFile + "` (result is FALSE)")
	} //end if
	if(smart.PathExists(theNewTestFile) != false) {
		fatalError("Failed to delete the test file `" + theNewTestFile + "` (path still exists after delete ...)")
	} //end if
	fmt.Println("Test PASSED ; DeleteFile: `" + theNewTestFile + "` ; Result =", isDSuccess)

	fmt.Println("------------------------- Copied File Delete TESTS -------------------------")

	isDCSuccess, errmsg := smart.SafePathFileDelete(theCopyTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while deleting the copied test file `" + theCopyTestFile + "`:", errmsg)
	} //end if
	if(isDCSuccess != true) {
		fatalError("Failed to delete the copied test file `" + theCopyTestFile + "` (result is FALSE)")
	} //end if
	if(smart.PathExists(theCopyTestFile) != false) {
		fatalError("Failed to delete the copied test file `" + theCopyTestFile + "` (path still exists after delete ...)")
	} //end if
	fmt.Println("Test PASSED ; DeleteCopiedFile: `" + theCopyTestFile + "` ; Result =", isDCSuccess)

	fmt.Println("------------------------- File Delete Non-Existing TESTS -------------------------")

	theTestFile = theTestFile + "-nonexisting"

	isDNSuccess, errmsg := smart.SafePathFileDelete(theTestFile, false)
	if(errmsg != "") {
		fatalError("Errors encountered while deleting (non-existing) the test file `" + theTestFile + "`:", errmsg)
	} //end if
	if(isDNSuccess != true) {
		fatalError("Failed to delete (non-existing) the test file `" + theTestFile + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; DeleteFile (non-existing): `" + theTestFile + "` ; Result =", isDNSuccess)

	fmt.Println("------------------------- Dir Delete TESTS -------------------------")

	isOkDirDeleteFinal, errmsg := smart.SafePathDirDelete(theBaseDir, false)
	if(errmsg != "") {
		fatalError("Errors encountered while deleting the test dir `" + theBaseDir + "`:", errmsg)
	} //end if
	if(isOkDirDeleteFinal != true) {
		fatalError("Failed to delete the test dir `" + theBaseDir + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; DirDelete: `" + theBaseDir + "` ; Result =", isOkDirDeleteFinal)

	if(smart.PathExists(theBaseDir)) {
		fatalError("Dir Tests EXISTS after delete: `" + theBaseDir)
	} //end if

	isOkDirDeleteNonExisting, errmsg := smart.SafePathDirDelete(theBaseDir, false)
	if(errmsg != "") {
		fatalError("Errors encountered while deleting the non-existing test dir `" + theBaseDir + "`:", errmsg)
	} //end if
	if(isOkDirDeleteNonExisting != true) {
		fatalError("Failed to delete the non-existing test dir `" + theBaseDir + "` (result is FALSE)")
	} //end if
	fmt.Println("Test PASSED ; DirNonExistingDelete: `" + theBaseDir + "` ; Result =", isOkDirDeleteNonExisting)

} //END FUNCTION


func main() {

	var input string = "Lorem Ipsum dolor sit Amet"

	LogToConsoleWithColors()

	fmt.Println("========================= DATE/TIME TESTS =========================")

	fmt.Println("---------- Date Now() ----------")
	testDTimeNow()
	fmt.Println("---------- Date Parse ----------")
	testDTimeParse()

	fmt.Println("========================= TRIM TESTS =========================")

	testAllTrimWhitespaces(input)

	fmt.Println("========================= HASH TESTS =========================")

	fmt.Println("MD5:", smart.Md5(input))
	fmt.Println("SHA1:", smart.Sha1(input))
	fmt.Println("SHA256:", smart.Sha256(input))
	fmt.Println("SHA384:", smart.Sha384(input))
	fmt.Println("SHA512:", smart.Sha512(input))

	fmt.Println("========================= BASE64 TESTS =========================")

	b64 := smart.Base64Encode(input)
	fmt.Println("B64-Enc:", b64)
	fmt.Println("B64-Dec:", smart.Base64Decode(b64))

	fmt.Println("========================= BIN/HEX TESTS =========================")

	hex := smart.Bin2Hex(input)
	fmt.Println("HEX-Enc:", hex)
	fmt.Println("HEX-Dec:", smart.Hex2Bin(hex))

	fmt.Println("========================= DATA ARCH/UNARCH TESTS =========================")

	arch := smart.DataArchive(input)
	fmt.Println("Data-Arch: `" + arch + "`")
	fmt.Println("Data-UnArch: `" + smart.DataUnArchive(arch) + "`")

	// INFO: arch data difers a little from PHP, maybe by some zlib metadata, but decrypt must work
	testPhpArchData := `HclBDkBAEETRw1hLplupZimDSMRKHMD06Psfgdj9/IfM1ZQ9Z00YLVlnfxNc+Zt+j6Phc+HM3tDkbcn7eR3tuU3SDKGhjwrCUaM4i6dbS7r9qRgEdIsq6i8=` + "\n" + `PHP.SF.151129/B64.ZLibRaw.HEX`;
	fmt.Println("Data-Arch-PHP: `" + testPhpArchData + "`")
	testPhpUnArchData := smart.DataUnArchive(testPhpArchData)
	fmt.Println("Data-UnArch-PHP: `" + testPhpUnArchData + "`")
	if(testPhpUnArchData != input) {
		fatalError("ERROR: DataArchive TEST Failed ... Archived Data is NOT EQUAL with Archived Data from PHP")
	} //end if

	fmt.Println("========================= BLOWFISH.CBC TESTS =========================")

	//--
	bfKey := "some.BlowFish! - Key@Test 2ks i782s982 s2hwgsjh2wsvng2wfs2w78s528 srt&^ # *&^&#*# e3hsfejwsfjh"
	//--
	bfInput := input + " " + smart.DateNowUtc()
	fmt.Println("Data-To-Encrypt: `" + bfInput + "`")
	bfData := smart.BlowfishEncryptCBC(bfInput, bfKey)
	fmt.Println("Data-Encrypted: `" + bfData + "`")
	testDecBfData := smart.BlowfishDecryptCBC(bfData, bfKey)
	fmt.Println("Data-Decrypted: `" + testDecBfData + "`")
	if((testDecBfData != bfInput) || (smart.Sha1(testDecBfData) != smart.Sha1(bfInput))) {
		fatalError("ERROR: BlowfishEncryptCBC TEST Failed ... Decrypted Data is NOT EQUAL with Plain Data")
	} //end if
	//--
	testPhpBfData := `695C491EF3E92DD8975423A91460F05F9DBBFDBE91DC55AE1D96CC43747B096D64CE08F42885D792505A56DF40CEE6B51FC399A3D756FADB4CE9A492BAE157B4B0DB0C6197D0E35B4C69F99266965686CB41628B75EA56CE006518F408CC0AF1`
	if(smart.BlowfishEncryptCBC(input, bfKey) != testPhpBfData) {
		fatalError("ERROR: BlowfishEncryptCBC TEST Failed ... Encrypted Data is NOT EQUAL with Encrypted Data from PHP")
	} //end if
	fmt.Println("Data-Encrypted-PHP: `" + testPhpBfData + "`")
	testDecPhpBfData := smart.BlowfishDecryptCBC(testPhpBfData, bfKey)
	fmt.Println("Data-Decrypted-PHP: `" + testDecPhpBfData + "`")
	if((testDecPhpBfData != input) || (smart.Sha1(testDecPhpBfData) != smart.Sha1(input))) {
		fatalError("ERROR: BlowfishDecryptCBC TEST Failed ... Decrypted Data is NOT EQUAL with Decrypted Data from PHP")
	} //end if

	fmt.Println("========================= STRING TESTS =========================")

//	var test2Str string = " Lorem Ipsum șȘțȚâÂăĂîÎ is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.       \r\n      \r\n\r\n"

	testStrings()

	fmt.Println("Strings SubStr Test Result: PASSED");

	fmt.Println("========================= UUID TESTS =========================")

	var u int = uid.UuidSessionSequence()
	u = uid.UuidSessionSequence()
	u = uid.UuidSessionSequence()
	if(u != 3) {
		fatalError("ERROR: UuidSessionSequence FAILED. Expected Result is 3 but get:", u)
	} //end if
	fmt.Println("Session Incremental UUID:", u)

	var uid10 = uid.Uuid1013Str(10)
	if((len(uid10) != 10) || (len(smart.StrTrimWhitespaces(uid10)) != 10)) {
		fatalError("ERROR: UuidStr(10) FAILED. Expected Result is a string of 10 characters length but get a string of :", len(uid10), "characters length, as: `" + uid10 + "`")
	} //end if
	fmt.Println("UUID-10:", uid10)

	var uid13 = uid.Uuid1013Str(13)
	if((len(uid13) != 13) || (len(smart.StrTrimWhitespaces(uid13)) != 13)) {
		fatalError("ERROR: UuidStr(13) FAILED. Expected Result is a string of 13 characters length but get a string of :", len(uid13), "characters length, as: `" + uid13 + "`")
	} //end if
	fmt.Println("UUID-13:", uid13)

	fmt.Println("========================= Json Encode TEST =========================")

	var theStrToEncAsJson string = "\n<c d=\"About:https://a.b?d=1&c=2\">aA-șȘțȚâÂăĂîÎ_+100.12345678901"
	fmt.Println("Json Encoded: `" + smart.JsonEncode(theStrToEncAsJson) + "` from `" + theStrToEncAsJson + "`")

	fmt.Println("========================= MarkerTPL TESTS =========================")

	var test3Str string = `{"a":2, "b":"\n<c d=\"About:https://a.b?d=1&c=2\">", "c":"aA-șȘțȚâÂăĂîÎ_+100.12345678901"}`
	var arr = map[string]string{
		"MARKER": 	test3Str,
		"MARKER2": 	`<Tyler="test"` + "\n" + `>`,
	}

	tpl := smart.MarkersTplRender(THE_TPL, arr, false, false)
	fmt.Println("\n" + "Raw TPL: `" + THE_TPL + "`" + "\n")
	fmt.Println("---------- ---------- ----------")
	fmt.Println("\n" + "Rendered TPL: `" + tpl + "`" + "\n")

	fmt.Println("========================= File System TESTS =========================")

	testFileSystem(input)

	fmt.Println("========================= # TESTS DONE # =========================")

} //END FUNCTION


// #END
