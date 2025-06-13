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

  signal status_array : slv_vector(8 downto 0)(16 downto 0) := (
    0 => "00000000000000000",
    1 => "00000000000000001",
    2 => "00000000000000010",
    3 => "00000000000000011",
    4 => "00000000000000100",
    5 => "00000000000000101",
    6 => "00000000000000110",
    7 => "00000000000000111",
    8 => "00000000000001000"
  );

  signal req : requester_out_t;
  signal com : completer_out_t;

begin

  clk <= not clk after 0.5 ns;

  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);

  afbd_main : entity afbd.main
  port map (
    clk_i => clk,
    rst_i => '0',

    apb_coms_i(0) => req,
    apb_coms_o(0) => com,

    status_array_i => status_array
  );

end architecture;
