library work;
  context work.cosim_context;
  use work.cosim.all;

library afbd;
  use afbd.main_pkg.all;


entity tb_cosim is
  generic(
    G_SW_GW_FIFO_PATH : string;
    G_GW_SW_FIFO_PATH : string
  );
end entity;


architecture test of tb_cosim is

  signal clk : std_logic := '0';

  signal cfg : std_logic_vector(31 downto 0);
  signal proc_out : my_proc_out_t;
  signal proc_in  : my_proc_in_t;
  signal exit_cnt : unsigned(3 downto 0) := (others => '0');

  signal req : requester_out_t;
  signal com : completer_out_t;

begin

  clk <= not clk after 0.5 ns;


  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);


  afbd_main : entity afbd.Main
  port map (
    clk_i => clk,
    rst_i => '0',

    apb_coms_i(0) => req,
    apb_coms_o(0) => com,

    Cfg_o      => cfg,
    my_Proc_o  => proc_out,
    my_proc_i  => proc_in,
    Exit_Cnt_i => std_logic_vector(exit_cnt)
  );


  proc_in.ret <= cfg;


  exit_cnt_driver : process (clk) is
  begin
    if rising_edge(clk) then
      if proc_out.exitt = '1' then
        exit_cnt <= exit_cnt + 1;
      end if;
    end if;
  end process;


  exitt_monitor : process (clk) is
    variable prev_exitt : std_logic := '0';
  begin
    if rising_edge(clk) then
      if prev_exitt = '1' and proc_out.exitt = '1' then
        report "proc_out.exitt asserted for more than one clock cycle"
          severity failure;
      end if;

      prev_exitt := proc_out.exitt;
    end if;
  end process;


  call_monitor : process (clk) is
  begin
    if rising_edge(clk) then
      assert proc_out.call = 'U' or proc_out.call = '0'
        report "proc_out.call shall not be asserted"
        severity failure;
    end if;
  end process;


end architecture;
