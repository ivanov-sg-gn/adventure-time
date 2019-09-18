<?php

class overload{
	public function get(){
		if(func_num_args() != 1){
			throw new Exception('Method expects 1 parameter');
		}
		
		$arg = func_get_arg(0);
		
		$typeParam = $this->getType($arg);

		switch($typeParam){
			case 'integer':
				$this->getById($arg);
				break;
			case 'array':
				$this->getByFilter($arg);
				break;
		}
	}	
	
	private function getType($param){
		$type = gettype($param);
		
		if($type == 'string'){
			if(DateTime::createFromFormat('Y-m-d H:i:s', $param) !== FALSE){
				return 'date';
			}
		}
		
		return gettype($param);
	}
	
	public function getById(int $id) : array{
		// ...
		
		return [];
	}
	
	public function getByFilter(array $arr) : array{
		// ...
		
		return [];
	}
	
}

$overload = new overload;

$overload->get(123); // run getById

$overload->get([]); // run getByFilter
