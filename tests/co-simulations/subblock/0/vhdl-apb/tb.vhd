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

  signal cfg0, cfg1, cfg2 : std_logic_vector(31 downto 0);

  signal req : requester_out_t;
  signal com : completer_out_t;

  signal blk0_req : requester_out_t;
  signal blk0_com : completer_out_t;

  signal blk1_req : requester_out_t;
  signal blk1_com : completer_out_t;

  signal blk2_req : requester_out_t;
  signal blk2_com : completer_out_t;

begin

  clk <= not clk after 0.5 ns;

  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);

  afbd_main : entity afbd.Main
  port map (
    clk_i => clk,
    rst_i => '0',
    reqs_i(0) => req,
    reqs_o(0) => com,
    blk0_coms_o(0) => blk0_req,
    blk0_coms_i(0) => blk0_com,
    blk1_coms_o(0) => blk1_req,
    blk1_coms_i(0) => blk1_com
  );


  afbd_blk0 : entity afbd.Blk0
  port map (
    clk_i => clk,
    rst_i => '0',
    reqs_i(0) => blk0_req,
    reqs_o(0) => blk0_com,
    cfg_o => cfg0,
    st_i => cfg0
  );


  afbd_blk1 : entity afbd.Blk1
  port map (
    clk_i => clk,
    rst_i => '0',
    reqs_i(0) => blk1_req,
    reqs_o(0) => blk1_com,
    blk2_coms_o(0) => blk2_req,
    blk2_coms_i(0) => blk2_com,
    cfg_o => cfg1,
    st_i => cfg1
  );


  afbd_blk2 : entity afbd.Blk2
  port map (
    clk_i => clk,
    rst_i => '0',
    reqs_i(0) => blk2_req,
    reqs_o(0) => blk2_com,
    cfg_o => cfg2,
    st_i => cfg2
  );

end architecture;
