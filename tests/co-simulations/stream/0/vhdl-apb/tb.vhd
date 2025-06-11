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

  signal fifo : fifo_t;
  signal fifo_stb : std_logic;
  signal val : unsigned(13 downto 0) := (others => '0');

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
    fifo_i     => fifo,
    fifo_stb_o => fifo_stb
  );


  fifo.val <= std_logic_vector(val);


  fifo_mock : process (clk) is
  begin
    if rising_edge(clk) then
      if fifo_stb = '1' then
        val <= val + 1;
      end if;
    end if;
  end process;

end architecture;
