
PUSH = 0
ADD = 1
MINUS = 2


def vmachine(program):
    program_counter = 0
    stack = list(range(100))
    stack_pointer = 0

    while program_counter < len(program):
        current_instruction = program[program_counter]
        if current_instruction == PUSH:
            stack[stack_pointer] = program[program_counter + 1]
            stack_pointer += 1
            program_counter += 1
        elif current_instruction == ADD:
            right = stack[stack_pointer - 1]
            stack_pointer -= 1
            left = stack[stack_pointer - 1]
            stack_pointer -= 1

            stack[stack_pointer] = left + right
            stack_pointer += 1
        elif current_instruction == MINUS:
            right = stack[stack_pointer - 1]
            stack_pointer-=1
            left = stack[stack_pointer - 1]
            stack_pointer -= 1

            stack[stack_pointer] = left - right
            stack_pointer += 1

        program_counter += 1

    print("stacktop: ", stack[stack_pointer - 1])


if __name__ == '__main__':
    program = [
        PUSH, 3,
        PUSH, 4,
        ADD,
        PUSH, 5,
        MINUS
    ]
    vmachine(program)

