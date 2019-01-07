<?php
// [@[#[!SF.DEV-ONLY!]#]@]
// Controller: Samples/BenchMark
// Route: ?/page/samples.benchmark (?page=samples.benchmark)
// Author: unix-world.org
// v.3.7.5 r.2018.03.09 / smart.framework.v.3.7

//----------------------------------------------------- PREVENT EXECUTION BEFORE RUNTIME READY
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------

define('SMART_APP_MODULE_AREA', 'INDEX'); // INDEX, ADMIN, SHARED

/**
 * Index Controller
 *
 * @ignore
 *
 */
class SmartAppIndexController extends SmartAbstractAppController {

	private $cart_currency = 'US$';

	public function Run() {

		//-- dissalow run this sample if not test mode enabled
		if(!defined('SMART_FRAMEWORK_TEST_MODE') OR (SMART_FRAMEWORK_TEST_MODE !== true)) {
			$this->PageViewSetErrorStatus(503, 'ERROR: Test mode is disabled ...');
			return;
		} //end if
		//--

		//--
		$op = $this->RequestVarGet('op', '', 'string');
		//--

		//--
		$cart = new \SmartModExtLib\Cart\ecommCart([
			'cartId' 			=> 'eCommCart',
			'cartMaxItem' 		=> 10, // Maximum item can added to cart, 0 = Unlimited
			'itemMaxQuantity' 	=> 50, // Maximum quantity of a item can be added to cart, 0 = Unlimited
			'useCookie' 		=> false // Do not use cookie, cart items will gone after browser closed
		]);
		//--
$yaml = <<<'YAML'
item:
	version: Roof Panels v.181017.1253
	id: 101
	type: p # s = Service | p = Product
	sku: ROOFPNLS
	name:
		en: Roof Pannels
		ro: Panouri Acoperiş
	currency: EUR
	price: 52.55
	discounts:
		10%
		12.5%
		15%
	tax: 19%
	um:
		unit: sqm
		hint:
			en: sq.m.
			ro: m.p.
	pak:
		umtype: 0 # U.M. Type: 0 = decimal ; 1 = integer
		umerg: 0 # U.M. Fractio: Ergonomic Order Qty.
		ummin: 0 # Minimum Order Qty.
	delivery:
		weight: # Shipping Weight in Kg
			unit: kg
			value: 5
		volume: # Shipping Volume in m3
			unit: m3
			value: 0.5
	warranty: # Warranty: -1 = lifetime ; 0 = no warranty ; 0+ = warranty (months, years, ...)
		unit: years
		value: 50
	image:
		source: "data:image/jpeg;base64,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"
		width: 200
		height: 200
	attributes:
		color:
			name:
				en: Color
				ro: Culoare
			type: list
			style:
				font-weight: bold
			validhint:
				en: Paint Color
				ro: Culoare vopsea
			validation:
				-
					id: RR29
					adjust-price: 1.28
					hint:
						en: Red
						ro: Roșu
					style:
						background-color: "#68200B"
						color: "#FFFFFF"
				-
					id: RR32
					hint: Dark Brown
					style:
						background-color: "#3F321E"
						color: "#FFFFFF"
				-
					id: RR11
					adjust-price: -0.50
					hint: Grey
		paint:
			name: Paint Mode
			type: free
			min: 0.25
			max: 0.5
			decimals: 2
			validation: decimal
			validhint: Paint thickness
		height:
			name:
				en: Panel Height (mm)
				ro: Înălțime Panou (mm)
			type: free
			min: 850
			max: 6000
			decimals: 0
			validation: integer
			validhint:
				en: Effective Length - Integer between
				ro: Lungime Efectivă - Număr întreg în intervalul
		width:
			name:
				en: Panel Width (mm)
				ro: Lățime Panou (mm)
			type: list
			style:
				width: 75px
			validhint:
				en: Effective Width - Integer selected
				ro: Lățime Efectivă - Număr întreg selectat
			validation:
				-
					id: 1100
				-
					id: 1200
		mountkit:
			name: Mount Kit
			optional: all
			type: free
			minlen: 5
			maxlen: 10
			validation: regex [a-z0-9\.]
			validhint: This is optional
		serial:
			name: Serial no.
			type: free
			optional: inventory
			length: 7
			validation: serial
			validhint: Required only for Inventory operations
	calculation:
		param: @ergonomic-qty@
		formula: ([[[@height]]] / 1000) * ([[[@width]]] / 1000) + (1 - 1)
		hint:
			en: Ergonomic Quantity Calculation by Width and Height
			ro: Calculare Cantitate Ergonomică in funcție de Lungime și Lățime
YAML;
$products = (new SmartYamlConverter())->parse($yaml);
/*
echo "<pre>";
echo Smart::escape_html(print_r($products,1));
echo "</pre>";
die();
*/
		//--
		if((string)$op == 'cart-json') {
			//--
			$this->PageViewSetCfg('rawpage', true);
			//--
			//print_r($_POST); die();
			$cart_op = $this->RequestVarGet('cart_action', '', 'string');
			$frm = $this->RequestVarGet('frm', [], 'array');
			//--
			$redirect = '';
			$message = '???';
			//-- Empty the cart
			if((string)$cart_op == 'empty') {
				$message = 'Cart cleared';
				$cart->clear();
				$cart->destroy();
				$redirect = '?page='.Smart::escape_url($this->ControllerGetParam('controller')).'&op=cart';
			} //end if
			//-- Add item
			if((string)$cart_op == 'add') {
				$cart_item_id = (string) $frm['id'];
				$cart_item_hash = (string) $frm['hash'];
				$cart_item_qty = (string) $frm['qty'];
				$cart_item_atts = (array) $frm['att'];
				$message = 'Product added';
				foreach($products as $key => $product) {
					if((string)$cart_item_id == (string)$product['id']) {
						break;
					} //end if
				} //end foreach
				$cart->add(
					[
						'currency' => (string) $this->cart_currency,
						'price' => (float) $product['price'],
						'tax' 	=> (int) $product['tax'],
						'atts' => (array) $this->parseItemAttributes($product['attributes'])
					],
					$product['id'],
					(array) $cart_item_atts,
					$cart_item_qty
				);
			} //end if
			//-- Update item
			if((string)$cart_op == 'update') {
				$redirect = '?page='.Smart::escape_url($this->ControllerGetParam('controller')).'&op=cart';
				if((string)$frm['cart'] == '@cart') {
					$message = 'Cart updated';
					$cart->multiupdate($frm);
				} else {
					$cart_item_id = (string) $frm['id'];
					$cart_item_hash = (string) $frm['hash'];
					$cart_item_qty = (string) $frm['qty'];
					$message = 'Product quantity updated';
					$cart->update(
						(string) $cart_item_id,
						(string) $cart_item_hash,
						$cart_item_qty
					);
				} //end if else
			} //end if
			//-- Remove item
			if((string)$cart_op == 'remove') {
				$cart_item_id = (string) $frm['id'];
				$cart_item_hash = (string) $frm['hash'];
			//	$cart_item_qty = (string) $frm['qty'];
				$message = 'Product removed';
				$redirect = '?page='.Smart::escape_url($this->ControllerGetParam('controller')).'&op=cart';
				foreach($products as $key => $product) {
					if((string)$cart_item_id == (string)$product['id']) {
						break;
					} //end if
				} //end foreach
				$cart->remove(
					$product['id'],
					(string) $cart_item_hash
				);
			} //end if
			//--
			$answer = 'OK';
			$title = 'Cart';
			//--
			$this->PageViewSetVar(
				'main',
				(string) SmartComponents::js_ajax_replyto_html_form($answer, $title, $message, $redirect)
			);
			return;
			//--
		} elseif((string)$op == 'cart') {
			//--
			$all_items = [];
			$cart_items = [];
			if(!$cart->isEmpty()) {
				$all_items = $cart->getItems();
				//print_r($all_items); die();
				foreach($all_items as $id => $items) {
					foreach($items as $key => $item) {
						foreach($products as $kk => $product) {
							if((string)$id == (string)$product['id']) {
								break;
							} //end if
						} //end if
//print_r($item); die();
						$tmp_arr = [];
						$tmp_arr['id'] = $item['id'];
						$tmp_arr['hash'] = $item['hash'];
						if((string)$product['qtytype'] == 'dec') {
							$tmp_arr['quantity'] = 0 + \Smart::format_number_dec($item['quantity'], 2, '.', '');
						} else {
							$tmp_arr['quantity'] = \Smart::format_number_int($item['quantity'], '+');
						} //end if else
						$tmp_arr['name'] = $this->parseItemName($product['name']);
						$tmp_arr['price'] = 0 + \Smart::format_number_dec($item['sell']['price'], 2, '.', '');
						$tmp_arr['tax'] = 0 + \Smart::format_number_dec($item['sell']['tax'], 2, '.', '');
						$tmp_arr['currency'] = $item['sell']['currency'];
						$tmp_arr['attributes'] = (array) $item['attributes'];
//print_r($item['attributes']); die();
						$cart_items[] = (array) $tmp_arr;
					} //end foreach
				} //end foreach
			} //end if
			//--
			$tpl = 'cart.mtpl.htm';
			$arr = [
				'PAGE-URL' 		=> (string) $this->ControllerGetParam('controller'),
				'CART-CURRENCY' 	=> (string) $this->cart_currency,
				'CART-TOTAL' 		=> (string) Smart::format_number_dec($cart->getCartTotal(), 2, '.', ''),
				'CART-ITEMS' 		=> (array) $cart_items
			];
			//--
		} elseif((string)$op == 'checkout') {
			//--
			return 404; // not implemented
			//--
		} else {
			//--
			$arr = [];
			if(is_array($products)) {
				foreach($products as $key => $val) {
					//print_r($products); die();
					//print_r($val); die();
					$attributes = (array) $this->parseItemAttributes($val['attributes']);
//print_r($attributes); die();
					$arr[] = [
						'id' 		=> (string) $val['id'],
						'name' 		=> (string) $this->parseItemName($val['name']),
						'price' 	=> (string) $val['price'],
						'currency' 	=> (string) $this->cart_currency,
						'img-src' 	=> (string) $val['image']['source'],
						'img-w' 	=> (string) $val['image']['width'],
						'img-h' 	=> (string) $val['image']['height'],
						'atts' 		=> (array)  $attributes,
						'hash' 		=> (string) $cart->calculateHash($val['id'], $attributes)
					];
				} //end foreach
			} //end if
			//--
			$tpl = 'shop.mtpl.htm';
			$arr = [
				'PAGE-URL' 		=> (string) $this->ControllerGetParam('controller'),
				'PRODUCTS-ARR' => (array) $arr
			];
			//--
		} //end if else
		//--
		$this->PageViewSetVars([
			'title' => 'eCommerce Test',
			'main' => SmartMarkersTemplating::render_file_template(
				(string) $this->ControllerGetParam('module-view-path').$tpl,
				(array) $arr,
				'no' // don't use caching (use of caching make sense only if file template is used more than once per execution)
			)
		]);
		//--

	} //END FUNCTION

	private function parseItemName($name) {
		//--
		return (string) (is_array($name) ? (string)$name['en'] : (string)$name);
		//--
	} //END FUNCTION


	private function parseItemAttributes($attribs) {
		//--
		$attributes = array();
		//--
		if(is_array($attribs)) {
			foreach($attribs as $k => $v) {
				$attvalues = array();
				if(is_array($v['name'])) {
					$attvalues['name'] = (array) $v['name'];
				} else {
					$attvalues['name'] = (string) $v['name'];
				} //end if else
				$attvalues['optional'] = (string) $v['optional'];
				$attvalues['type'] = (string) $v['type'];
				if(((string)$attvalues['type'] == 'list') AND (\Smart::array_size($v['validation']) > 0)) {
					$attvalues['validation'] = array();
					foreach($v['validation'] as $kk => $vv) {
						if(is_array($vv)) {
							$vv['id'] = (string) trim((string)$vv['id']);
							if((string)$vv['id'] != '') {
								$attvalues['validation'][(string)$vv['id']] = [
									'adjust-price' => 0 + \Smart::format_number_dec($vv['adjust-price'], 2, '.', ''),
									'hint' => (\Smart::array_size($vv['hint']) > 0) ? (array) $vv['hint'] : (string) $vv['hint'],
									'style' => (\Smart::array_size($vv['style']) > 0) ? (array) $vv['style'] : (string) $vv['style'],
								];
							} //end if
						} //end if
					} //end foreach
				} else {
					$attvalues['validation'] = (string) $v['validation'];
					if((string)$attvalues['validation'] == 'integer') {
						$attvalues['min'] = (int) \Smart::format_number_int($v['min'], '+');
						$attvalues['max'] = (int) \Smart::format_number_dec($v['max'], '+');
					} elseif((string)$attvalues['validation'] == 'decimal') {
						$attvalues['decimals'] = (int) (($v['decimals'] > 0 AND $v['decimals'] <= 4) ? $v['decimals'] : 0);
						$attvalues['min'] = (string) \Smart::format_number_dec($v['min'], (int)$v['decimals'], '.', '');
						$attvalues['max'] = (string) \Smart::format_number_dec($v['max'], (int)$v['decimals'], '.', '');
					} else { // string
						$attvalues['length'] = (int) ($v['length'] > 0 ? $v['length'] : 0);
						if(!$attvalues['length']) {
							$attvalues['minlen'] = (int) ($v['minlen'] > 0 ? $v['minlen'] : 0);
							$attvalues['maxlen'] = (int) ($v['maxlen'] > 0 ? $v['maxlen'] : 0);
						} //end if
					} //end if else
				} //end if
				if(is_array($v['validhint'])) {
					$attvalues['validhint'] = (array) $v['validhint'];
				} else {
					$attvalues['validhint'] = (string) $v['validhint'];
				} //end if else
				if(is_array($v['style'])) {
					$attvalues['style'] = (array) $v['style'];
				} else {
					$attvalues['style'] = (string) $v['style'];
				} //end if else
				$attributes[(string)$k] = (array) $attvalues;
			} //end foreach
		} //end if
		//--
		//print_r($attributes); die();
		return (array) $attributes;
		//--
	} //END FUNCTION


} //END CLASS

//end of php code
?>