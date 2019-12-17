package body def_monitor is
   protected body monitor is
      entry sentarse when sillas > 0 is
      begin
         sillas := sillas - 1;
         esperandoComida := esperandoComida + 1;
      end sillasLock;

      procedure levantarse is
      begin
         sillas := sillas + 1;
      end sillasUnlock;

      entry comer when comidaPreparada > 0 is
      begin
         comidaPreparada := comidaPreparada - 1;
         esperandoComida := esperandoComida - 1;
      end comer;

      procedure darComida is
      begin
         comidaPreparada := comidaPreparada + 1;
      end cocinar;


   end monitor;
end def_monitor;
