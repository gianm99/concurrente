package def_monitor is

   protected type monitor is
      entry sentarse;
      procedure levantarse;
      entry comer;
      procedure darComida;
      procedure dormir;
      function esperando return integer;
      function dormidos return integer;
      function sillas_libres return integer;
   private
      sillas: integer := 4;
      enanos_esperando: integer := 0;
      enanos_dormidos: integer := 0;
      comida: integer := 0;
   end monitor;

end def_monitor;
