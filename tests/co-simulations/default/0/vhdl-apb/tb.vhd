library work;
  context work.cosim_context;
  use work.cosim.all;


entity tb_cosim is
  generic(
    G_SW_GW_FIFO_PATH : string;
    G_GW_SW_FIFO_PATH : string
  );
end entity;


architecture test of tb_cosim is

  signal clk : std_logic := '0';

  signal cfg0 : std_logic_vector(5 downto 0);
  signal cfg1 : std_logic_vector(3 downto 0);

  signal req : requester_out_t;
  signal com : completer_out_t;

begin

  clk <= not clk after 0.5 ns;

  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);

  afbd_main : entity afbd.Main
  port map (
    clk_i => clk,
    rst_i => '0',
    reqs_i(0) => req,
    reqs_o(0) => com,
    cfg0_o => cfg0,
    st0_i => cfg0,
    cfg1_o => cfg1,
    st1_i => cfg1
  );

end architecture;
