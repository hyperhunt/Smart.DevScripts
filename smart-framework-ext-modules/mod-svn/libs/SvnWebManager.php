<?php
// [LIB - SmartFramework / Svn / Svn Web Manager]
// (c) 2006-2017 unix-world.org - all rights reserved

// Class: \SmartModExtLib\Svn\SvnWebManager
// Type: Module Library

namespace SmartModExtLib\Svn;

//----------------------------------------------------- PREVENT DIRECT EXECUTION
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------


// # Sample Configuration #
/*
//-- SVN related configuration (add this in etc/config-admin.php)
$configs['svn']['cmd'] = '/usr/local/bin/svn';
$configs['svn']['7za'] = '/usr/local/bin/7za';
$configs['svn']['repos'] = [
	'netvision-enterprise' => [ 'url' => 'https://repos/svn/some/', 'user' => 'user123', 'pass' => 'pass123' ]
];
//--
*/

//=====================================================================================
//===================================================================================== CLASS START
//=====================================================================================


final class SvnWebManager {

	// ::
	// v.170917

	private static $svn_cache_dir = 'tmp/cache/svn/'; 		// must have trailing slash :: the svn proc jail root


	//============================================================ OK
	public static function listRepos() {
		//--
		$repos = (array) \Smart::get_from_config('svn.repos');
		if(\Smart::array_size($repos) <= 0) {
			return array();
		} //end if
		//--
		$arr = array();
		//--
		foreach($repos as $key => $val) {
			//--
			if(((string)trim((string)$key) != '') AND (self::validateCfgRepoEntry($val))) {
				//--
				$tmp_arr = (array) self::execSvnCmd('info', (string)$val['url'], '', $val['user'], $val['pass'], 'xml-arr'); // OK
				//print_r($tmp_arr); die();
				//--
				if(\Smart::array_size($tmp_arr) > 0) {
					if(\Smart::array_size($tmp_arr['info']) > 0) {
						if(\Smart::array_size($tmp_arr['info'][0]) > 0) {
							//--
							$arr[] = [
								'repo-name' 		=> (string) trim((string)$key),
								'last-rev-num' 		=> (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'info.0.entry.0.commit|@attributes.0.revision', '.'),
								'last-rev-author' 	=> (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'info.0.entry.0.commit.0.author.0', '.'),
								'last-rev-date' 	=> (string) date('D, d M Y H:i:s', strtotime((string)\Smart::array_get_by_key_path((array)$tmp_arr, 'info.0.entry.0.commit.0.date.0', '.')))
							];
							//print_r($arr); die();
							//--
						} //end if
					} //end if
				} //end if
				//--
				$tmp_arr = array();
				//--
			} //end if
			//--
		} //end foreach
		//--
		return (array) $arr;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	public static function listRepo($repo, $path, $rev) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return array();
		} //end if
		//--
		$arr = array();
		//--
		$tmp_arr = (array) self::execSvnCmd('list', (string)$rdata['url'], (string)$path, $rdata['user'], $rdata['pass'], 'xml-arr', [ 'rev' => (string)$rev ]); // OK
		//print_r($tmp_arr); die();
		if(\Smart::array_size($tmp_arr) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN List: Invalid XML (1): ['.$path.']',
				'ERR: SVN List: Invalid XML (1)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		//--
		if(\Smart::array_size($tmp_arr['lists']) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN List: Invalid XML (2): ['.$path.']',
				'ERR: SVN List: Invalid XML (2)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		//--
		if(\Smart::array_size($tmp_arr['lists'][0]) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN List: Invalid XML (3): ['.$path.']',
				'ERR: SVN List: Invalid XML (3)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		//--
		$tmp_arr = (array) $tmp_arr['lists'][0];
		//print_r($tmp_arr); die();
		//--
		if((\Smart::array_size($tmp_arr['list']) <= 0) OR (\Smart::array_size($tmp_arr['list|@attributes']) <= 0)) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN List: Invalid XML (4): ['.$path.']',
				'ERR: SVN List: Invalid XML (4)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		//--
		$tmp_path = (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'list|@attributes.0.path', '.');
		$tmp_entratt = \Smart::array_get_by_key_path((array)$tmp_arr, 'list.0.entry|@attributes', '.'); // don't force array !! because (array)'' = array(0=>'')
		$tmp_entries = \Smart::array_get_by_key_path((array)$tmp_arr, 'list.0.entry', '.'); // don't force array !!
		$tmp_arr = array();
		if(\Smart::array_size($tmp_entries) <= 0) {
			//var_dump((array)$tmp_entries); die((string)phpversion());
			return array(); // Fix: no entries found
		} //end if
		//print_r($tmp_entries); die();
		if(\Smart::array_size($tmp_entratt) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN List: Invalid XML (5): ['.$path.']',
				'ERR: SVN List: Invalid XML (5)' // msg to display
			);
			die(''); // just in case
			return array(); // no entries found
		} //end if
		//--
		for($i=0; $i<\Smart::array_size($tmp_entries); $i++) {
			//--
			$val = (array) $tmp_entries[$i];
			//print_r($val); die();
			//-- proc
			if(\Smart::array_size($val) > 0) {
				//--
				$size = '-';
				if(array_key_exists('size', (array)$val)) {
					$size = (string) \SmartUtils::pretty_print_bytes((int)\Smart::array_get_by_key_path((array)$val, 'size.0', '.'), 1, ' ');
				} //end if
				//--
				$the_item_name = (string) \Smart::array_get_by_key_path((array)$val, 'name.0', '.');
				$the_item_type = (string) \Smart::array_get_by_key_path((array)$tmp_entratt, $i.'.kind', '.');
				$the_item_icon_suffix = '';
				if((string)$the_item_type == 'file') {
					$the_item_icon_suffix = (string) self::getFileTypeIcon($the_item_name);
					if((string)$the_item_icon_suffix != '') {
						$the_item_icon_suffix = '-'.$the_item_icon_suffix;
					} //end if
				} //end if
				//--
				$arr[] = [
					'repo-name' => (string) $repo,
					'base-path' => (string) $tmp_path,
					'type' => (string) $the_item_type,
					'name' => (string) $the_item_name,
					'icon-suffix' => (string) $the_item_icon_suffix,
					'size' => (string) $size,
					'last-rev-num' => (string) \Smart::array_get_by_key_path((array)$val, 'commit|@attributes.0.revision', '.'),
					'last-rev-author' => (string) \Smart::array_get_by_key_path((array)$val, 'commit.0.author.0', '.'),
					'last-rev-date' => (string) date('D, d M Y H:i:s', strtotime((string)\Smart::array_get_by_key_path((array)$val, 'commit.0.date.0', '.'))),
				];
				//--
				//print_r($arr); die();
				//--
			} //end if
			//--
		} //end for
		//--
		//print_r($arr); die();
		//--
		return (array) $arr;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	public static function getHeadRevision($repo) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return -1;
		} //end if
		//--
		$tmp_arr = (array) self::execSvnCmd('get-revs', (string)$rdata['url'], '', $rdata['user'], $rdata['pass'], 'xml-arr', [ 'start-rev' => 'HEAD', 'num-revs' => 1 ]); // get latest revision
		//print_r($tmp_arr); die();
		$entry_zero = (array) \Smart::array_get_by_key_path((array)$tmp_arr, 'log.0.logentry|@attributes.0', '.');
		if((\Smart::array_size($entry_zero) <= 0) OR (!array_key_exists('revision', (array)$entry_zero))) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Export: Failed to get SVN Export Head Revision ...',
				'ERR: Failed to get SVN Export Head Revision' // msg to display
			);
			die(''); // just in case
		} //end if
		//print_r($entry_zero); die();
		//--
		return (string) ($entry_zero['revision'] ? (int)$entry_zero['revision'] : '');
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	public static function getCompare($repo, $path, $rev) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return array();
		} //end if
		//-- get compare just for root and later filter by path
		$tmp_arr = self::execSvnCmd('compare', (string)$rdata['url'], '/', $rdata['user'], $rdata['pass'], 'xml-arr', [ 'rev' => (string)$rev ]); // OK
		//print_r($tmp_arr); die();
		//--
		if(\Smart::array_size($tmp_arr['log']) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Compare: Invalid XML (1): ['.$path.']',
				'ERR: SVN Compare: Invalid XML (1)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		if(\Smart::array_size($tmp_arr['log'][0]) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Compare: Invalid XML (2): ['.$path.']',
				'ERR: SVN Compare: Invalid XML (2)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		//--
		$tmp_arr = (array) $tmp_arr['log'][0];
		$tmp_rev = (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'logentry|@attributes.0.revision', '.');
		$tmp_author = (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'logentry.0.author.0', '.');
		$tmp_date = (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'logentry.0.date.0', '.');
		$tmp_msg = (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'logentry.0.msg.0', '.');
		$tmp_arr = (array) \Smart::array_get_by_key_path((array)$tmp_arr, 'logentry.0.paths.0', '.');
		//--
		if((string)$tmp_rev !== (string)$rev) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Compare: Invalid XML Rev.['.$tmp_rev.'] / Rev.['.$rev.'] (3): ['.$path.']',
				'ERR: SVN Compare: Invalid XML (3)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		if((\Smart::array_size($tmp_arr['path']) <= 0) OR (\Smart::array_size($tmp_arr['path|@attributes']) <= 0) OR (\Smart::array_size($tmp_arr['path']) != \Smart::array_size($tmp_arr['path|@attributes']))) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Compare: Invalid XML (4): ['.$path.']',
				'ERR: SVN Compare: Invalid XML (4)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		//--
		$arr = array();
		$arr['metainfo'] = [
			'rev' 		=> (string) $tmp_rev,
			'author' 	=> (string) $tmp_author,
			'date' 		=> (string) date('D, d M Y H:i:s', strtotime((string)$tmp_date)),
			'msg' 		=> (string) $tmp_msg
		];
		$arr['changes'] = [];
		for($i=0; $i<\Smart::array_size($tmp_arr['path']); $i++) {
			if(((string)$path == '') OR ((string)$path == '/') OR (strpos($tmp_arr['path'][$i], $path) === 0)) {
				$tmp_atts = (array) $tmp_arr['path|@attributes'][$i];
				$the_item_icon_suffix = '';
				if((string)$tmp_atts['kind'] == 'file') {
					$the_item_icon_suffix = (string) self::getFileTypeIcon((string) $tmp_arr['path'][$i]);
					if((string)$the_item_icon_suffix != '') {
						$the_item_icon_suffix = '-'.$the_item_icon_suffix;
					} //end if
				} //end if
				$arr['changes'][] = [
					'path' 			=> (string) $tmp_arr['path'][$i],
					'type' 			=> (string) $tmp_atts['kind'],
					'icon-suffix' 	=> (string) $the_item_icon_suffix,
					'txtmod' 		=> (string) $tmp_atts['text-mods'],
					'action' 		=> (string) $tmp_atts['action'],
					'propsmod' 		=> (string) $tmp_atts['prop-mods'],
					'copyfromp' 	=> (string) $tmp_atts['copyfrom-path'],
					'copyfromr' 	=> (string) $tmp_atts['copyfrom-rev']
				];
			} //end if
		} //end for
		//--
		//print_r($arr); die();
		return (array) $arr;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	public static function getFile($repo, $path, $rev) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return '';
		} //end if
		//--
		return (string) self::execSvnCmd('cat', (string)$rdata['url'], (string)$path, $rdata['user'], $rdata['pass'], 'string', [ 'rev' => (string)$rev ]); // OK
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ aaa
	public static function getRealPathFromPrevRevision($repo, $path, $revx, $revy) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return '';
		} //end if
		//--
		$tmp_arr = (array) self::execSvnCmd('prev-info', (string)$rdata['url'], (string)$path, $rdata['user'], $rdata['pass'], 'xml-arr', [ 'rev-old' => (string)$revx, 'rev-new' => (string)$revy ]); // OK
		//print_r($tmp_arr); die();
		$path = (string) \Smart::array_get_by_key_path((array)$tmp_arr, 'info.0.entry.0.relative-url.0', '.');
		$rpath = (string) trim((string)substr((string)$path, 1));
		if((substr((string)$path, 0, 1) != '^') OR ((string)$rpath == '')) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Revision Path: Invalid Real Path for: ['.$path.'] @ ['.$path.']',
				'ERR: Invalid SVN Revision Path' // msg to display
			);
			die(''); // just in case
		} //end if
		//die('Path: '.$rpath);
		//--
		return (string) $rpath;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	public static function getDiffFile($repo, $path, $revx, $revy) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return '';
		} //end if
		//--
		return (string) self::execSvnCmd('diff', (string)$rdata['url'], (string)$path, $rdata['user'], $rdata['pass'], 'string', [ 'rev-old' => (string)$revx, 'rev-new' => (string)$revy ]); // OK
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	public static function exportPath($repo, $path, $rev) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return array();
		} //end if
		//--
		$path = (string) trim((string)$path);
		if(((string)$path == '') OR ((string)$path == '/')) {
			$finpath = '';
		} else {
			$finpath = (string) '_'.trim((string)\Smart::safe_filename((string)$path, '-'), '-');
		} //end if
		$expname = (string) \Smart::safe_filename((string)$repo.'.r'.$rev.$finpath);
		$expdir = (string) 'svn-exp/'.\Smart::uuid_10_seq().'-'.\Smart::uuid_10_str().'-'.\Smart::uuid_10_num();
		$archdir = (string) \SmartFileSysUtils::add_dir_last_slash((string)$expdir).$expname;
		if(!\SmartFileSysUtils::check_if_safe_path($archdir)) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Export: Invalid Export Dir:['.$archdir.']',
				'ERR: Invalid SVN Export Dir' // msg to display
			);
			die(''); // just in case
		} //end if
		//--
		\SmartFileSystem::dir_create((string)self::$svn_cache_dir.$expdir, true); // recursive
		//--
		$ok = self::execSvnCmd('export', (string)$rdata['url'], (string)$path, $rdata['user'], $rdata['pass'], 'exit-code', [ 'export-dir' => (string)$archdir, 'rev' => (string)$rev ]); // OK
		$archname = '';
		$fcontent = '';
		if(($ok) AND (\SmartFileSystem::is_type_dir((string)self::$svn_cache_dir.$archdir))) {
			$archname = (string) self::buildArchive((string)$expdir, (string)$expname, '7z');
			if((string)$archname != '') {
				$fcontent = (string) \SmartFileSystem::read((string)$archname);
			} //end if
		} //end if
		//--
		\SmartFileSystem::dir_delete((string)self::$svn_cache_dir.$expdir, true);
		//--
		return array(
			'f-content' => (string) $fcontent,
			'f-mime' => (string) ($archname ? 'application/zip' : ''),
			'f-name' => (string) ($archname ? (string) \SmartFileSysUtils::get_file_name_from_path((string)$archname) : '')
		);
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	public static function listRevs($repo, $path, $start, $num) {
		//--
		$repo = (string) trim((string)$repo);
		$repos = (array) \Smart::get_from_config('svn.repos');
		$rdata = (array) $repos[(string)trim((string)$repo)];
		if(((string)trim((string)$repo) == '') OR (!self::validateCfgRepoEntry($rdata))) {
			return array();
		} //end if
		//--
		$tmp_arr = (array) self::execSvnCmd('get-revs', (string)$rdata['url'], (string)$path, $rdata['user'], $rdata['pass'], 'xml-arr', [ 'start-rev' => (string)$start, 'num-revs' => (int)$num, 'rev' => (string)$start ]);
		//print_r($tmp_arr); die();
		//--
		if(\Smart::array_size($tmp_arr['log']) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Rev. List: Invalid XML (1): ['.$path.']',
				'ERR: SVN Rev. List: Invalid XML (1)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		if(\Smart::array_size($tmp_arr['log'][0]) <= 0) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Rev. List: Invalid XML (2): ['.$path.']',
				'ERR: SVN Rev. List: Invalid XML (2)' // msg to display
			);
			die(''); // just in case
			return array();
		} //end if
		//--
		$entries = (array) $tmp_arr['log'][0];
		$tmp_arr = array();
		$arr = [];
		//print_r($entries); die();
		for($i=0; $i<\Smart::array_size($entries['logentry']); $i++) {
			//--
			$arr[] = [
				'revision' => (string) \Smart::array_get_by_key_path((array)$entries['logentry|@attributes'], $i.'.revision', '.'),
				'author' => (string) \Smart::array_get_by_key_path((array)$entries['logentry'], $i.'.author.0', '.'),
				'date' => (string) date('D, d M Y H:i:s', strtotime((string)\Smart::array_get_by_key_path((array)$entries['logentry'], $i.'.date.0', '.'))),
				'msg' => (string) \Smart::array_get_by_key_path((array)$entries['logentry'], $i.'.msg.0', '.')
			];
			//print_r($arr); die();
			//--
		} //end for
		//print_r($arr); die();
		//--
		return (array) $arr;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================
	public static function isTextFileByMimeType($mimetype) {
		//--
		if(in_array((string)$mimetype, [
			'application/x-php',
			'application/javascript',
			'application/json',
			'application/xml',
			'text/html',
			'text/css',
			'text/plain',
			'text/x-vcard',
			'text/calendar',
			'text/csv',
			'image/svg+xml',
			'text/ldif',
			'application/pgp-signature'
		])) {
			$out = true;
		} else {
			$out = false;
		} //end if else
		//--
		return (bool) $out;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================
	public static function getFileTypeIcon($path) {
		//--
		$file = (string) \SmartFileSysUtils::get_file_name_from_path($path);
		//--
		if(in_array((string)strtolower((string)$file), [
			'makefile'
		])) {
			return 'shell';
		} //end if
		//--
		if(in_array((string)strtolower((string)$file), [
			'license',
			'readme'
		])) {
			return 'txt';
		} //end if
		//--
		$ext = (string) \SmartFileSysUtils::get_file_extension_from_path($path);
		$type = '';
		switch((string)strtolower((string)$ext)) { // SYNC WITH \SmartFileSysUtils::mime_eval()
			case 'cs': // C#
			case 'c': // C
			case 'y': // Yacc source code file
			case 'cpp': // C++
			case 'ypp': // Bison source code file
			case 'cxx': // C++
			case 'yxx': // Bison source code file
			case 'm': // Objective-C Method
				$type = 'c';
				break;
			case 'h': // C header
			case 'hpp': // C++ header
			case 'hxx': // C++ header
				$type = 'h';
				break;
			case 'txt': // text
			case 'log': // log file
				$type = 'txt';
				break;
			case 'xhtml':
			case 'xml':
			case 'xsl':
			case 'dtd':
				$type = 'xml';
				break;
			case 'htm':
			case 'html':
				$type = 'html';
				break;
			case 'css':
			case 'less':
			case 'scss':
			case 'sass':
				$type = 'css';
				break;
			case 'js':
				$type = 'js';
				break;
			case 'php':
				$type = 'php';
				break;
			case 'pl':
				$type = 'pl';
				break;
			case 'py':
				$type = 'py';
				break;
			case 'java':
				$type = 'java';
				break;
			case 'sql':
				$type = 'db';
				break;
			case 'csh': // C-Shell script
			case 'sh': // shell script
			case 'awk': // AWK script
			case 'cmd': // windows command file
			case 'bat': // windows batch file
				$type = 'shell';
				break;
			case 'svg':
				$type = 'svg';
				break;
			case 'png':
			case 'gif':
			case 'jpg':
			case 'jpe':
			case 'jpeg':
			case 'tif':
			case 'tiff':
			case 'wmf':
			case 'bmp':
				$type = 'photo';
				break;
			default:
				$type = ''; // not recognized or icon n/a
		} //end if
		//--
		return (string) $type;
		//--
	} //END FUNCTION
	//============================================================


	//##### PRIVATES


	//============================================================ OK
	private static function validateCfgRepoEntry($repo_entry) {
		//--
		if(!\Smart::array_size($repo_entry) > 0) {
			return false;
		} //end if
		//--
		if((string)trim((string)$repo_entry['url']) == '') {
			return false;
		} //end if
		//--
		if($repo_entry['user'] !== null) {
			if((string)trim((string)$repo_entry['user']) == '') {
				return false;
			} //end if
		} //end if
		//--
		if($repo_entry['pass'] !== null) {
			if((string)trim((string)$repo_entry['pass']) == '') {
				return false;
			} //end if
		} //end if
		//--
		return true;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	private static function execSvnCmd($what, $repo, $path, $auth_user, $auh_pass, $format, $options=[]) {
		//--
		$cmd = (string) self::buildSvnCmd($what, $repo, $path, $auth_user, $auh_pass, $options);
		if((string)$cmd == '') {
			return array();
		} //end if
		//--
		$exearr = (array) \SmartUtils::run_proc_cmd((string)$cmd, null, (string)self::$svn_cache_dir); // avoid proc open in web root !!
		if(($exearr['exitcode'] !== 0) OR ((string)$exearr['stderr'] != '')) {
			\Smart::raise_error(
				__METHOD__.' #ERR# SVN Command:['.$cmd.'] Returned Some Errors ; ExitCode=['.$exearr['exitcode'].'] ; ErrorMsg: '.$exearr['stderr'],
				'ERR: Errors when running command' // msg to display
			);
			die(''); // just in case
		} //end if
		$out = (string) trim((string)$exearr['stdout']);
		$exearr = array(); // free mem
		//--
		switch((string)$format) {
			case 'xml-arr':
				if((string)$out == '') {
					\Smart::raise_error(
						__METHOD__.' #ERR# SVN Command:['.$cmd.'] Returned Empty Output ...',
						'ERR: Errors when running command' // msg to display
					);
					die(''); // just in case
				} //end if
				//die((string)$out);
				$arr = (array) (new \SmartXmlParser('extended'))->transform((string)$out);
				if(\Smart::array_size($arr) <= 0) {
					\Smart::raise_error(
						__METHOD__.' #ERR# SVN Command:['.$cmd.'] Returned Invalid Output:['."\n".$out."\n".']',
						'ERR: Errors when running command' // msg to display
					);
					die(''); // just in case
				} //end if
				return (array) $arr;
				break;
			case 'string':
				return (string) $out;
				break;
			case 'exit-code':
				return true; // return TRUE as the real exit code was checked above
			default:
				if(\Smart::array_size($arr) <= 0) {
					\Smart::raise_error(
						__METHOD__.' #ERR# SVN Command:['.$cmd.'] Invalid command Output Type selected:['.$format.']',
						'ERR: Invalid command Output Type selected' // msg to display
					);
					die(''); // just in case
				} //end if
		} //end switch
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	private static function buildSvnCmd($what, $repo, $path, $auth_user, $auh_pass, $options) {
		//--
		$repo = (string) trim((string)$repo);
		if(strpos((string)$repo, '://') === false) {
			$repo = 'file:///INVALID-REPOSITORY/--invalid-repo-name-err-svn--';
		} //end if
		if(substr((string)$repo, -1, 1) == '/') {
			$repo = (string) substr((string)$repo, 0, -1);
		} //end if
		//--
		$path = (string) trim((string)$path);
		if(substr((string)$path, 0, 1) != '/') {
			$path = (string) '/'.$path;
		} //end if
		//--
		//die($repo.$path);
		//--
		$cmdsvn = (string) trim((string)\Smart::get_from_config('svn.cmd'));
		if((string)$cmdsvn == '') {
			return '';
		} //end if
		//--
		$base_cmd = (string) self::escapeCmdExe((string)$cmdsvn).' --non-interactive';
		$base_cmd .= ' --config-dir '.self::escapeCmdArg('svn-cfg'); // this path is relative to the proc jailroot
		if($auth_user !== null) {
			$base_cmd .= ' --no-auth-cache --username '.self::escapeCmdArg((string)$auth_user);
			if($auh_pass !== null) {
				$base_cmd .= ' --password '.self::escapeCmdArg((string)$auh_pass);
			} //end if
		} //end if
		//--
		$cmd = '';
		switch((string)$what) {
			case 'info':
				$cmd = (string) $base_cmd.' --xml info '.self::escapeCmdArg((string)$repo);
				break;
			case 'prev-info': // get real path (maybe changed from previous revisions)
				$cmd = (string) $base_cmd.' --xml info --revision '.self::escapeCmdArg((string)$options['rev-old']).' '.self::escapeCmdArg((string)$repo.$path.'@'.(string)$options['rev-new']);
				break;
			case 'list':
				if((string)trim((string)$options['rev']) != '') {
					$rev = (string) $options['rev'];
				} else {
					$rev = 'HEAD';
				} //end if else
				$cmd = (string) $base_cmd.' --xml ls --revision '.self::escapeCmdArg((string)$rev).' '.self::escapeCmdArg((string)$repo.$path.'@'.(string)$rev);
				break;
			case 'cat':
				if((string)trim((string)$options['rev']) != '') {
					$rev = (string) $options['rev'];
				} else {
					$rev = 'HEAD';
				} //end if else
				$cmd = (string) $base_cmd.' cat --revision '.self::escapeCmdArg((string)$rev).' '.self::escapeCmdArg((string)$repo.$path.'@'.(string)$rev);
				break;
			case 'compare': // compare changes between revisions
				if((string)trim((string)$options['rev']) != '') {
					$rev = (string) $options['rev'];
				} else {
					$rev = 'HEAD';
				} //end if else
				$cmd = (string) $base_cmd.' --xml log --verbose --revision '.self::escapeCmdArg((string)$options['rev']).':'.self::escapeCmdArg((string)$options['rev']).' '.self::escapeCmdArg((string)$repo.$path.'@'.(string)$options['rev']);
				//die($cmd);
				break;
			case 'diff': // show diff on a text file
				$cmd = (string) $base_cmd.' diff --revision '.self::escapeCmdArg((string)$options['rev-old']).':'.self::escapeCmdArg((string)$options['rev-new']).' '.self::escapeCmdArg((string)$repo.$path.'@');
				break;
			case 'export':
				if((string)$options['export-dir'] == '') {
					\Smart::raise_error(
						__METHOD__.' #ERR# SVN Export: Empty or Invalid Export Dir:['.$options['export-dir'].']',
						'ERR: Invalid Dir for SVN Export' // msg to display
					);
					die(''); // just in case
				} //end if
				if((string)trim((string)$options['rev']) != '') {
					$rev = (string) $options['rev'];
				} else {
					$rev = 'HEAD';
				} //end if else
				$cmd = (string) $base_cmd.' export --revision '.self::escapeCmdArg((string)$rev).' '.self::escapeCmdArg((string)$repo.$path.'@'.(string)$rev).' '.self::escapeCmdArg((string)$options['export-dir'].'@');
				break;
			case 'get-revs':
				$start_rev = (string) $options['start-rev'];
				$num_revs = (int) $options['num-revs'];
				if((string)trim((string)$options['rev']) != '') {
					$rev = (string) $options['rev'];
				} else { // $rev is required for removed or changed paths
					$rev = 'HEAD';
				} //end if else
				if($num_revs <= 0) {
					$num_revs = 1;
				} //end if
				$cmd = (string) $base_cmd.' --xml log --revision '.self::escapeCmdArg($start_rev).':0 --limit '.(int)$num_revs.' '.self::escapeCmdArg((string)$repo.$path.'@'.(string)$rev);
				break;
			case 'get-head':
				$cmd = (string) $base_cmd.' --xml info '.self::escapeCmdArg((string)$repo.$path.'@'.(string)$rev);
				break;
			default:
				// nothing
		} //end switch
		//--
		//die($cmd);
		return (string) $cmd;
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	private static function buildArchive($dir, $archname, $format) {
		//--
		$dir = (string) \Smart::safe_pathname((string)\SmartFileSysUtils::add_dir_last_slash((string)$dir));
		\SmartFileSysUtils::raise_error_if_unsafe_path($dir);
		$arch_dir = (string) \SmartFileSysUtils::add_dir_last_slash(\Smart::safe_filename((string)$archname)); // archive name dir
		\SmartFileSysUtils::raise_error_if_unsafe_path($arch_dir);
		$arch_file = (string) \Smart::safe_filename((string)$archname.'.7z'); // archive name
		\SmartFileSysUtils::raise_error_if_unsafe_path($arch_file);
		$fpatharch = (string) \Smart::safe_pathname((string)self::$svn_cache_dir.$dir.$arch_file);
		\SmartFileSysUtils::raise_error_if_unsafe_path($fpatharch);
		//--
		$cmd7za = (string) trim((string)\Smart::get_from_config('svn.7za'));
		if((string)$cmd7za == '') {
			\Smart::raise_error(
				__METHOD__.' #ERR# Empty Archive Command: 7za',
				'ERR: Empty Archive Command' // msg to display
			);
			die(''); // just in case
		} //end if
		//--
		$cmd = (string) self::escapeCmdExe((string)$cmd7za).' a -ssc -t7z -m0=lzma '.self::escapeCmdArg((string)$arch_file).' '.self::escapeCmdArg((string)$arch_dir);
		$exearr = (array) \SmartUtils::run_proc_cmd((string)$cmd, null, (string)self::$svn_cache_dir.$dir);
		if(($exearr['exitcode'] !== 0) OR ((string)$exearr['stderr'] != '')) {
			\Smart::raise_error(
				__METHOD__.' #ERR# Archive Command:['.$cmd.'] Returned Some Errors ; ExitCode=['.$exearr['exitcode'].'] ; ErrorMsg: '.$exearr['stderr'],
				'ERR: Errors when running archive command' // msg to display
			);
			die(''); // just in case
		} //end if
		//die(trim((string)$exearr['stdout'])); // not needed for archiving
		$exearr = array(); // free mem
		//--
		if(\SmartFileSystem::is_type_file((string)$fpatharch)) {
			if(filesize((string)$fpatharch) > 0) {
				//--
				$cmd = (string) self::escapeCmdExe((string)$cmd7za).' t '.self::escapeCmdArg((string)$arch_file);
				$exearr = (array) \SmartUtils::run_proc_cmd((string)$cmd, null, (string)self::$svn_cache_dir.$dir);
				//die(trim((string)$exearr['stdout'])); // not needed for archiving
				if(($exearr['exitcode'] !== 0) OR ((string)$exearr['stderr'] != '') OR (stripos((string)$exearr['stdout'], "\n".'Everything is Ok'."\n") === false)) {
					\Smart::raise_error(
						__METHOD__.' #ERR# Archive Test Command:['.$cmd.'] Returned Some Errors ; ExitCode=['.$exearr['exitcode'].'] ; ErrorMsg: '.$exearr['stderr'].' ; StdOut: '.$exearr['stdout'],
						'ERR: Errors when running archive test command' // msg to display
					);
					die(''); // just in case
				} //end if
				$exearr = array(); // free mem
				//--
				return (string) $fpatharch;
				//--
			} //end if
		} //end if
		//--
		return '';
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	private static function escapeCmdExe($cmd) {
		//--
		$cmd = (string) trim((string)\Smart::normalize_spaces((string)$cmd));
		//--
		return (string) escapeshellcmd((string)$cmd);
		//--
	} //END FUNCTION
	//============================================================


	//============================================================ OK
	private static function escapeCmdArg($arg) {
		//--
		$arg = (string) trim((string)\Smart::normalize_spaces((string)$arg));
		//--
		return (string) escapeshellarg((string)$arg);
		//--
	} //END FUNCTION
	//============================================================


} //END CLASS


//=====================================================================================
//===================================================================================== CLASS END
//=====================================================================================


// end of php code
?>