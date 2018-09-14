<?php
// Controller: Agile, FlowChart
// Route: ?page=agile.flowchart
// Author: Radu I.
// r.2018-08-28

//----------------------------------------------------- PREVENT EXECUTION BEFORE RUNTIME READY
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------

define('SMART_APP_MODULE_AREA', 'ADMIN');
define('SMART_APP_MODULE_AUTH', true);

class SmartAppAdminController extends SmartAbstractAppController {


	public function Run() {

		//--
		$sq_rd = (array) (new \SmartModDataModel\Agile\SqFlowcharts())->getAllByUuid();
		//--
		$this->PageViewSetVars([
			'title' 	=> 'Agile :: FlowCharts / List',
			'main' 			=> SmartMarkersTemplating::render_file_template(
				$this->ControllerGetParam('module-path').'views/flowchart.htm', // the view
				[
					'FLOW-CHARTS' => (array) $sq_rd
				]
			)
		]);
		//--



	} // END FUNCTION


} // END CLASS

?>