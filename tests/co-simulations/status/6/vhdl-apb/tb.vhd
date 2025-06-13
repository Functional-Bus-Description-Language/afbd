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

  signal sts : slv_vector(2 downto 0)(33 downto 0);

  signal req : requester_out_t;
  signal com : completer_out_t;

begin

   clk <= not clk after 0.5 ns;

  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);

  sts(0) <= std_logic_vector(S0(33 downto 0));
  sts(1) <= std_logic_vector(S1(33 downto 0));
  sts(2) <= std_logic_vector(S2(33 downto 0));


  afbd_main : entity afbd.Main
  port map (
    clk_i => clk,
    rst_i => '0',

    apb_coms_i(0) => req,
    apb_coms_o(0) => com,

    sts_i => sts
  );

end architecture;
