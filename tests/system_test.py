import unittest
import subprocess
import hashlib

PATH_TO_EXECUTABLE = 'build/godd.bin'


class SystemTest(unittest.TestCase):
    def setUp(self) -> None:
        SystemTest.generate_test_file_1()
        SystemTest.generate_test_file_2()
        SystemTest.generate_test_file_3()

    @staticmethod
    def generate_test_file_1():
        with open('test_file_1.raw', mode='wb') as output:
            output.write(b'Hello World')

    @staticmethod
    def generate_test_file_2():
        with open('test_file_2.raw', mode='wb') as output:
            for _ in range(50):
                output.write(b'Hello World')

    @staticmethod
    def generate_test_file_3():
        with open('test_file_3.raw', mode='wb') as output:
            for _ in range(500):
                output.write(b'Hello World')

    def test_file_to_file_1(self):
        self.copy_file_and_compare_hash(input_file='test_file_1.raw', output_file='test_output_1.out')

    def test_file_to_file_2(self):
        self.copy_file_and_compare_hash(input_file='test_file_2.raw', output_file='test_output_2.out')

    def test_file_to_file_3(self):
        self.copy_file_and_compare_hash(input_file='test_file_3.raw', output_file='test_output_3.out')

    def copy_file_and_compare_hash(self, input_file: str, output_file: str):
        result = subprocess.run([PATH_TO_EXECUTABLE] + ['-if', input_file, '-of', output_file])
        if result.returncode != 0:
            self.fail(result.stderr)

        self.assertEqual(SystemTest.get_file_size(input_file), SystemTest.get_file_size(output_file))

        hash_original = SystemTest.calculate_file_hash(input_file)
        hash_new_file = SystemTest.calculate_file_hash(output_file)
        self.assertEqual(hash_original, hash_new_file)

    @staticmethod
    def calculate_file_hash(file_path: str) -> str:
        hash_object = hashlib.new('sha1')
        hash_object.update(SystemTest.get_file_content(file_path=file_path))
        return hash_object.hexdigest()

    @staticmethod
    def get_file_size(file_path: str) -> int:
        return len(SystemTest.get_file_content(file_path=file_path))

    @staticmethod
    def get_file_content(file_path: str) -> bytes:
        with open(file_path, 'rb') as file:
            return file.read()

    def test_standard_input_1(self):
        self.copy_from_standard_input(input_file='test_file_1.raw', output_file='test_output_stdin_1.out')

    def test_standard_input_2(self):
        self.copy_from_standard_input(input_file='test_file_2.raw', output_file='test_output_stdin_2.out')

    def test_standard_input_3(self):
        self.copy_from_standard_input(input_file='test_file_2.raw', output_file='test_output_stdin_3.out')

    def copy_from_standard_input(self, input_file: str, output_file: str):
        result = subprocess.run(f"cat {input_file} | {PATH_TO_EXECUTABLE} -of {output_file}", shell=True)
        self.assertEqual(result.returncode, 0)

        self.assertEqual(SystemTest.get_file_size(input_file), SystemTest.get_file_size(output_file))

        hash_original = SystemTest.calculate_file_hash(input_file)
        hash_new_file = SystemTest.calculate_file_hash(output_file)
        self.assertEqual(hash_original, hash_new_file)

    def test_standard_output_1(self):
        self.copy_to_standard_output(input_file='test_file_1.raw')

    def test_standard_output_2(self):
        self.copy_to_standard_output(input_file='test_file_2.raw')

    def test_standard_output_3(self):
        self.copy_to_standard_output(input_file='test_file_2.raw')

    def copy_to_standard_output(self, input_file: str):
        result = subprocess.run(f"{PATH_TO_EXECUTABLE} -if {input_file}", shell=True, capture_output=True)
        self.assertEqual(result.returncode, 0)

        self.assertEqual(SystemTest.get_file_size(input_file), len(result.stdout))

        hash_original = SystemTest.calculate_file_hash(input_file)

        hash_object = hashlib.new('sha1')
        hash_object.update(result.stdout)

        self.assertEqual(hash_original, hash_object.hexdigest())
