package def_monitor is

   protected type monitor is
      entry sentarse;
      procedure levantarse;
      entry comer;
      procedure darComida;
   private
      sillas: integer := 4;
      esperandoComida: integer := 0;
      comidaPreparada: integer := 0;
   end monitor;

end def_monitor;
