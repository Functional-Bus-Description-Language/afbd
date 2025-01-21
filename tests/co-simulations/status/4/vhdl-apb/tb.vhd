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

  signal st : unsigned(39 downto 0) := (others => '1');

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
    st_i => std_logic_vector(st)
  );

  counter : process (clk) is
  begin
    if rising_edge(clk) then
      st(39 downto 32) <= st(39 downto 32) - 1;
    end if;
  end process;

end architecture;
