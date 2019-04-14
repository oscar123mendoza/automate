module AutomateCtl
  class Ctl
    attr_accessor :command_map, :fh_output, :verbose
    attr_reader :exe_name

    def initialize()
      @fh_output = STDOUT
      @exe_name = File.basename($0)
      @command_map = {
        "help" => {
          :arity => 1,
          :desc => "Print this help message."
        }
      }
    end

    def self.to_method_name(name)
      name.gsub(/-/, '_').to_sym
    end

    def to_method_name(name)
      Ctl.to_method_name(name)
    end

    def load_files(path)
      Dir["#{path}/*.rb"].each do |file|
        load_file(file)
      end
    end

    def load_file(filepath)
      eval(IO.read(filepath), nil, filepath, 1)
    end

    def add_command(name, description, arity = 1, &block)
      @command_map[name] = { desc: description, arity: arity }
      self.class.send(:define_method, to_method_name(name).to_sym) { |*args| block.call(*args) }
    end

    def log(msg)
      fh_output.puts msg
    end

    def run_command(command)
      system(command)
      $?
    end

    def help(*args)
      log "#{exe_name}: command (subcommand)\n"
      command_map.keys.sort.each do |command|
        log command
        log "  #{command_map[command][:desc]}"
      end
    end

    # If it begins with a '-', it is an option.
    def is_option?(arg)
      arg && arg[0] == '-'
    end

    def help_requested?(args)
      args.include?("-h") || args.include?("--help") || args.include?("help")
    end

    def run(args)
      cmd, *cmd_args = args[0]
      if cmd.nil?
        help
        exit 1
      elsif help_requested?([cmd])
        help
        exit 0
      end

      # returns either hash content of command or nil
      command = command_map[cmd]
      if command.nil?
        log "I don't know that command."
        help
        exit 1
      end

      if args.length > 1 && command[:arity] != 2
        log "The command #{command_to_run} does not accept any arguments"
        exit 2
      end

      File::umask(022)
      run_root_check!
      method_to_call = to_method_name(cmd)
      send(method_to_call, *cmd_args)
    end

    def run_root_check!
      if root_required? && Process.uid != 0
        $stderr.puts "You must run this command as root. Try adding sudo to the start of your command."
        exit 1
      end
    end

    def root_required?
      ! (inside_docker_container? || help_requested?(ARGV))
    end

    def inside_docker_container?
      # Builds will hang unless we short-circuit
      return false if ENV['TRAVIS'] == 'true'

      file_path = "/proc/1/cgroup"
      if File.exist?(file_path)
        File.read(file_path).include?("docker")
      else
        false
      end
    end
  end
end
