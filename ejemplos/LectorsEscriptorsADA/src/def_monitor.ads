package def_monitor is

   protected type RWMonitor is
      entry readerLock;
      procedure readerUnlock;
      entry writerLock;
      procedure writerUnlock;
   private
      readers : integer := 0;
      writing : boolean := false;
   end RWMonitor;

end def_monitor;
