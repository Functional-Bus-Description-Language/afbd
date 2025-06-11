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

  signal foo : foo_out_t;
  signal count : std_logic_vector(31 downto 0) := (others => '0');

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
    foo_o   => foo,
    count_i => count
  );


  adder : process (clk) is
  begin
    if rising_edge(clk) then
      if foo.call = '1' then
        count <= std_logic_vector(unsigned(count) + to_unsigned(1, count'length));
      end if;
    end if;
  end process;

end architecture;
