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

  signal st : std_logic_vector(to_integer(apb.main_pkg.C) - 1 downto 0) := std_logic_vector(to_unsigned(to_integer(apb.main_pkg.C), to_integer(apb.main_pkg.C)));
  signal stl : slv_vector(1 downto 0)(7 downto 0) := (
    0 => std_logic_vector(to_unsigned(to_integer(apb.main_pkg.CL(0)), 8)),
    1 => std_logic_vector(to_unsigned(to_integer(apb.main_pkg.CL(1)), 8))
  );

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
    st_i => st,
    stl_i => stl
  );

end architecture;
