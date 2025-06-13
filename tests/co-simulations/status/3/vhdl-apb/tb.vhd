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

  signal status_array0 : slv_vector(7 downto 0)(7 downto 0) := (
    0 => "00000000",
    1 => "00000001",
    2 => "00000010",
    3 => "00000011",
    4 => "00000100",
    5 => "00000101",
    6 => "00000110",
    7 => "00000111"
  );

  signal status_array1 : slv_vector(3 downto 0)(4 downto 0) := (
    0 => "00000",
    1 => "00001",
    2 => "00010",
    3 => "00011"
  );

  signal status_array2 : slv_vector(5 downto 0)(6 downto 0) := (
    0 => "0000000",
    1 => "0000001",
    2 => "0000010",
    3 => "0000011",
    4 => "0000100",
    5 => "0000101"
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

    apb_coms_i(0) => req,
    apb_coms_o(0) => com,

    status_array0_i => status_array0,
    status_array1_i => status_array1,
    status_array2_i => status_array2
  );

end architecture;
