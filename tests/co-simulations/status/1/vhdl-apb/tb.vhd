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

  signal lower : std_logic_vector(29 downto 0);
  signal upper : std_logic_vector(19 downto 0);
  signal st    : std_logic_vector(49 downto 0);

  signal req : requester_out_t;
  signal com : completer_out_t;

begin

  clk <= not clk after 0.5 ns;

  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);

  afbd_main : entity afbd.Main
  port map (
    clk_i => clk,
    rst_i => '0',
    coms_i(0) => req,
    coms_o(0) => com,
    lower_o => lower,
    upper_o => upper,
    st_i => st
  );

  st <= upper & lower;

end architecture;
