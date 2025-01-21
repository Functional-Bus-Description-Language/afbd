library ieee;
  use ieee.std_logic_1164.all;
  use ieee.numeric_std.all;

library std;
  use std.textio.all;

library lapb;
  use lapb.apb.all;
  use lapb.bfm;


package cosim is

  procedure cosim_interface (
    constant read_fifo_path : string;
    constant write_fifo_path : string;
    signal clk : in std_logic;
    signal req : inout requester_out_t;
    signal com : in  completer_out_t;
    constant bfm_config : bfm.config_t := bfm.DEFAULT_CONFIG
 );

end package;


package body cosim is

  procedure cosim_interface (
    constant read_fifo_path : string;
    constant write_fifo_path : string;
    signal clk : in std_logic;
    signal req : inout requester_out_t;
    signal com : in  completer_out_t;
    constant bfm_config : bfm.config_t := bfm.DEFAULT_CONFIG
  ) is
    file wr_pipe : text;
    file rd_pipe : text;

    variable code    : character;
    variable rd_line : line;
    variable wr_line : line;

    variable addr : unsigned(31 downto 0);
    variable data : std_logic_vector(31 downto 0);

    variable end_status : integer;
  begin
    file_open(rd_pipe, read_fifo_path, read_mode);
    file_open(wr_pipe, write_fifo_path, write_mode);

    while not endfile(rd_pipe) loop
      readline(rd_pipe, rd_line);
      read(rd_line, code);

      case code is
      when 'W' =>
        hread(rd_line, addr);
        read(rd_line, code);
        if code /= ',' then
          report "wrong separator in the write command" severity error;
        end if;
        hread(rd_line, data);

        bfm.write(addr, data, clk, req, com, cfg => bfm_config);

        write(wr_line, string'("ACK"));
        writeline(wr_pipe, wr_line);
        flush(wr_pipe);
      when 'R' =>
        hread(rd_line, addr);

        bfm.read(addr, clk, req, com, cfg => bfm_config);

        write(wr_line, to_string(com.rdata));
        writeline(wr_pipe, wr_line);
        flush(wr_pipe);
      when 'T' =>
        hread(rd_line, data);

        wait for to_integer(unsigned(data)) * 1 ns;

        write(wr_line, string'("ACK"));
        writeline(wr_pipe, wr_line);
        flush(wr_pipe);
      when 'E' =>
        read(rd_line, end_status);

        write(wr_line, string'("ACK"));
        writeline(wr_pipe, wr_line);
        flush(wr_pipe);

        if end_status /= 0 then
          report "end status " & integer'image(end_status) & ", check proper log in /tmp/afbd/..." severity failure;
        end if;

        file_close(rd_pipe);
        file_close(wr_pipe);
        std.env.finish;
      when others =>
        report "cosim interface - unknown command: '" & code & "'" severity failure;
      end case;
    end loop;
  end procedure;

end package body;
