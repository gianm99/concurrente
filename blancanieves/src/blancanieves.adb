-- Autor: Gian Lucas Martin Chamorro
-- Enlace: https://youtu.be/7UQQAyvB-dE
with Ada.Text_IO,Ada.Integer_Text_IO,def_monitor;
use Ada.Text_IO,Ada.Integer_Text_IO,def_monitor;

procedure blancanieves is
   -- Enumerado de los nombres de los enanos
   type nombres_enanos is(SABIO,GRUNON,FELIZ,DORMILON,TIMIDO,MOCOSO,TONTIN);
   -- Constante de veces que come un enano
   MAX: constant integer:=2;
   -----------------------------
   --  Tipo protegido para la SC
   -----------------------------
   m: monitor;
   -----------------------------
   --  Especificacion de procedimiento(s)
   -----------------------------
   procedure mostrar_estado;
   -----------------------------
   --  Cuerpo de procedimiento(s)
   -----------------------------
   procedure mostrar_estado is
   begin
      -- Muestra la cantidad de enanos esperando para comer y las sillas libres
      Put_Line("esperando = " & m.esperando'img & " sillas = " & m.sillas_libres'img);
   end mostrar_estado;
   -----------------------------
   --  Especificaciones de las tareas
   -----------------------------
   task type tarea_enano is
      entry Start(nombre: in nombres_enanos);
   end tarea_enano;

   task type tarea_blancanieves is
      entry Start;
   end tarea_blancanieves;

   -----------------------------
   -- Cuerpos de las tareas
   -----------------------------
   task body tarea_enano is
      mi_nombre : nombres_enanos;
   begin
      -- Esperar a que se inicie la tarea y asignar el nombre
      accept Start (nombre: in nombres_enanos) do
         mi_nombre:=nombre;
      end Start;
      -- Mientras le queden repeticiones por hacer
      for I in 1..MAX loop
         Put_Line(mi_nombre'img & " va a trabajar a la mina");
         delay Duration(4.0);
         m.sentarse;
         Put_Line(mi_nombre'img & " se sienta");
         m.comer;
         mostrar_estado;
         Put_Line("-----------------> "& mi_nombre'img & " come!!!");
         delay Duration(1.5);
         m.levantarse;
         mostrar_estado;
      end loop;
      m.dormir;
      Put_Line(mi_nombre'img & " se va a DORMIR " & m.dormidos'img & "/7");
   end tarea_enano;

   task body tarea_blancanieves is
   begin
      -- Esperar a que se inicie la tarea
      accept Start;
      -- Mientras queden enanos despiertos
      while m.dormidos<7 loop
         -- Mientras queden enanos esperando comida
         while m.esperando>0 loop
            -- Cocina para un enano
            Put_Line("BLANCANIEVES cocina para un enano");
            delay Duration(0.5);
            m.darComida;
         end loop;
         -- Se va a pasear cuando no hay enanos esperando para comer
         Put_Line("BLANCANIEVES se va a pasear");
         delay Duration(1.5);
      end loop;
      -- Cuando estan dormidos los enanos, se va a dormir y acaba
      Put_Line("BLANCANIEVES se va a DORMIR");
   end tarea_blancanieves;

   -----------------------------
   --  Tareas
   -----------------------------
   type Blancanieves is new tarea_blancanieves;
   b: Blancanieves;
   type Enanos is array (nombres_enanos) of tarea_enano;
   e:Enanos;
begin
   -- Empezar todas las tareas
   b.Start; -- Blancanieves
   for I in e'Range loop
      e(I).Start(I); -- Enanos
   end loop;
end blancanieves;
