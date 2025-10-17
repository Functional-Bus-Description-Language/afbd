library work;
  context work.cosim_context;
  use work.cosim.all;

library ltypes;
  use ltypes.types.all;


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

  signal blk_req : requester_out_array_t(0 to 2);
  signal blk_com : completer_out_array_t(0 to 2);

  signal x : slv_vector(0 to 2)(31 downto 0);
  signal sum : std_logic_vector(33 downto 0);

begin

  clk <= not clk after 0.5 ns;

  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);

  afbd_main : entity afbd.Main
  port map (
    clk_i => clk,
    rst_i => '0',

    apb_coms_i(0) => req,
    apb_coms_o(0) => com,

    blk_apb_reqs_o(0) => blk_req(0),
    blk_apb_reqs_o(1) => blk_req(1),
    blk_apb_reqs_o(2) => blk_req(2),
    blk_apb_reqs_i(0) => blk_com(0),
    blk_apb_reqs_i(1) => blk_com(1),
    blk_apb_reqs_i(2) => blk_com(2),

    sum_i => sum
  );

  gen_blks : for i in 0 to 2 generate
    blk : entity afbd.Blk
    port map (
      clk_i => clk,
      rst_i => '0',

      apb_coms_i(0) => blk_req(i),
      apb_coms_o(0) => blk_com(i),

      X_o => x(i)
    );
  end generate;


  adder : process (clk) is
  begin
    if rising_edge(clk) then
      sum <= std_logic_vector(
        resize(unsigned(x(0)), sum'length) +
        resize(unsigned(x(1)), sum'length) +
        resize(unsigned(x(2)), sum'length)
      );
    end if;
  end process;

end architecture;
