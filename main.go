package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// initProject 函数用于初始化项目
func initProject(projectName string) error {
	// 获取当前工作目录
	var content string
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 拼接项目目录路径
	projectPath := currentDir

	// 创建项目目录
	//err = os.MkdirAll(projectPath, 0755)
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("项目目录 %s 创建成功。\n", projectPath)

	// 创建 README.md 文件
	readmePath := filepath.Join(projectPath, "README.md")
	readmeFile, err := os.Create(readmePath)
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	// 向 README.md 文件写入内容
	writer := bufio.NewWriter(readmeFile)
	_, err = writer.WriteString(fmt.Sprintf("# %s 项目说明\n", projectPath))
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	fmt.Println("README.md 文件创建成功。")
	// CMakeLists.txt
	cmakeListPath := filepath.Join(projectPath, "CMakeLists.txt")
	cmakeListFile, err := os.Create(cmakeListPath)
	if err != nil {
		return err
	}
	defer cmakeListFile.Close()
	writer = bufio.NewWriter(cmakeListFile)
	content = fmt.Sprintf(`cmake_minimum_required(VERSION 3.15...3.27)
project(%s) # Replace 'my_project' with the name of your project

if (CMAKE_VERSION VERSION_LESS 3.18)
  set(DEV_MODULE Development)
else()
  set(DEV_MODULE Development.Module)
endif()

if (WIN32) 
  set(findPython "where.exe")
else()
  set(findPython "which")
endif()

execute_process(
  COMMAND "${findPython}" python 
  OUTPUT_STRIP_TRAILING_WHITESPACE OUTPUT_VARIABLE Python_EXECUTABLE)
# set(Python3_ROOT_DIR "/Users/thw/miniforge3/envs/pycpp")
# set(Python_EXECUTABLE "/Users/thw/miniforge3/envs/pycpp/bin/python")
find_package(Python 3.8 COMPONENTS Interpreter ${DEV_MODULE} REQUIRED)

if (NOT CMAKE_BUILD_TYPE AND NOT CMAKE_CONFIGURATION_TYPES)
  set(CMAKE_BUILD_TYPE Release CACHE STRING "Choose the type of build." FORCE)
  set_property(CACHE CMAKE_BUILD_TYPE PROPERTY STRINGS "Debug" "Release" "MinSizeRel" "RelWithDebInfo")
endif()

# Detect the installed nanobind package and import it into CMake
execute_process(
  COMMAND "${Python_EXECUTABLE}" -m nanobind --cmake_dir
  OUTPUT_STRIP_TRAILING_WHITESPACE OUTPUT_VARIABLE nanobind_ROOT)
find_package(nanobind CONFIG REQUIRED)
nanobind_add_module(%s %s.cpp)
`, projectName, projectName, projectName)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	fmt.Println("cmake 文件创建成功。")
	// 创建 main.go 文件
	newfile := fmt.Sprintf("%s.cpp", projectName)
	mainGoPath := filepath.Join(projectPath, newfile)
	mainGoFile, err := os.Create(mainGoPath)
	if err != nil {
		return err
	}
	defer mainGoFile.Close()

	// 向 my_ext 文件写入内容
	writer = bufio.NewWriter(mainGoFile)
	content = fmt.Sprintf(`#include <nanobind/nanobind.h>
namespace nb = nanobind;
int add(int a, int b) { return a + b; }

class Rectangle {
    private:
        // 矩形的长和宽
        double length;
        double width;
    
    public:
        // 构造函数，用于初始化矩形的长和宽
        Rectangle(double l, double w) {
            length = l;
            width = w;
        }
    
        // 计算矩形的面积
        double area() {
            return length * width;
        }
    
        // 计算矩形的周长
        double perimeter() {
            return 2 * (length + width);
        }
    };
    
int fib(int n) {
    if (n <= 1) return n;
    return fib(n - 1) + fib(n - 2);
}

NB_MODULE(%s, m) {
    m.def("add", &add);
    m.def("fib", &fib);
    nb::class_<Rectangle>(m, "Rectangle")
        .def(nb::init<double, double>())
        .def("area", &Rectangle::area)
        .def("perimeter", &Rectangle::perimeter);
}
`, projectName)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	fmt.Println("cpp 文件创建成功。")

	return nil
}

func main() {
	if len(os.Args) != 3 || os.Args[1] != "init" {
		fmt.Printf("用法: %s init <项目名称>\n", os.Args[0])
		os.Exit(1)
	}

	projectName := os.Args[2]
	err := initProject(projectName)
	if err != nil {
		fmt.Printf("初始化项目时出错: %v\n", err)
		os.Exit(1)
	}
}
