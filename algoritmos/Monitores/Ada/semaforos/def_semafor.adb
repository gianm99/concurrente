protected body Semafor is
	entry Wait when Contador>0 is
	begin
		Contador:=Contador-1;
	end Wait;

	procedure Signal is
	begin
		Contador:=Contador+1;
	end Signal;
end Semafor;