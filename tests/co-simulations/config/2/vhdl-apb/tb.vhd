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

  signal cfg : std_logic_vector(47 downto 0);
  constant ZERO : std_logic_vector(47 downto 0) := (others => '0');

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
    cfg_o => cfg
  );


  atomicity_guardian : process (clk) is
  begin
    if rising_edge(clk) then
      assert cfg = ZERO or cfg = std_logic_vector(resize(apb.main_pkg.VALID_VALUE, 48))
        report "invalid value: " & to_hstring(cfg) & ", expecting: " & to_hstring(std_logic_vector(apb.main_pkg.VALID_VALUE))
        severity failure;
    end if;
  end process;

end architecture;
