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

  signal add : add_out_t;
  signal result : std_logic_vector(31 downto 0) := (others => '0');

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

    add_o    => add,
    result_i => result
  );


  adder : process (clk) is
  begin
    if rising_edge(clk) then
      if add.call = '1' then
        result <= std_logic_vector(
          resize(unsigned(add.terms(0)), result'length) +
          resize(unsigned(add.terms(1)), result'length) +
          resize(unsigned(add.terms(2)), result'length) +
          resize(unsigned(add.terms(3)), result'length) +
          resize(unsigned(add.terms(4)), result'length) +
          resize(unsigned(add.terms(5)), result'length) +
          resize(unsigned(add.terms(6)), result'length) +
          resize(unsigned(add.terms(7)), result'length) +
          resize(unsigned(add.terms(8)), result'length) +
          resize(unsigned(add.terms(9)), result'length)
        );
      end if;
    end if;
  end process;

end architecture;
