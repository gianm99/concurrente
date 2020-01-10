n: Integer := 0;
pragma Volatile(n);
S: Semafor(1);

task type Tasca_comptador;

task body Tasca_comptador is
begin
	for i in 1..1000000 loop
		S.wait;
		n:=n+1;
		S.Signal;
	end loop;
end Tasca_comptador;