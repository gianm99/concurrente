with Text_IO,def_monitor,Ada.Strings.Unbounded;
use Text_IO,def_monitor,Ada.Strings.Unbounded;

procedure Main is

   ENANOS : constant integer := 7;
   type enanos is(Doc,Grumpy,Happy,Sleepy,Bashful,Sneezy,Dopey);

   --  Tipo protegido para la SC
   -----------------------------
   m: monitor;
   --  Especificaciones  de las tareas
   task type Enano is
      entry Start (
   end Enano;

begin
   --  Insert code here.
   null;
end Main;
