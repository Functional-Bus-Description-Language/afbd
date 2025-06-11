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

  signal req : requester_out_t;
  signal com : completer_out_t;

  -- Testbench specific signals
  signal add : add_t;
  signal add_stb : std_logic;

  signal result : result_t;
  signal result_stb : std_logic;

  signal buff : slv_vector(0 to to_integer(apb.main_pkg.DEPTH) - 1)(40 downto 0);
  signal buff_write_ptr, buff_read_ptr : natural := 0;

begin

  clk <= not clk after 0.5 ns;

  cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, req, com);

  afbd_main : entity afbd.Main
  port map (
    clk_i => clk,
    rst_i => '0',
    coms_i(0) => req,
    coms_o(0) => com,

    add_o     => add,
    add_stb_o => add_stb,

    result_i     => result,
    result_stb_o => result_stb
  );


  Write_Driver : process (clk) is
  begin
    if rising_edge(clk) then
      if add_stb = '1' then
        buff(buff_write_ptr) <= std_logic_vector(
          resize(unsigned(add.a), buff(0)'length) +
          resize(unsigned(add.b), buff(0)'length) +
          resize(unsigned(add.c), buff(0)'length)
        );
        buff_write_ptr <= buff_write_ptr + 1;
      end if;
    end if;
  end process;


  result.res <= buff(buff_read_ptr);


  Read_Driver : process (clk) is
  begin
    if rising_edge(clk) then
      if result_stb = '1' then
        buff_read_ptr <= (buff_read_ptr + 1) mod to_integer(apb.main_pkg.DEPTH);
      end if;
    end if;
  end process;

end architecture;
