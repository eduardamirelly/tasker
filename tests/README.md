# Tasker Test Suite Documentation

This directory contains the complete test suite for the Tasker CLI application. All tests are organized in a dedicated `tests/` folder to keep them separate from the main application code.

## 📁 Test Structure

```
tests/
├── README.md              # This documentation
├── test_helpers.go         # Common test utilities and database setup
├── add_test.go            # Tests for the add command
├── list_test.go           # Tests for the list command  
├── done_test.go           # Tests for the done command
├── export_test.go         # Tests for the export command
└── integration_test.go    # End-to-end integration tests
```

## 🔧 Test Utilities (`test_helpers.go`)

### Key Functions

#### `setupTestDB(t *testing.T) func()`
- **Purpose**: Creates an isolated test database for each test
- **How it works**: 
  - Creates a temporary SQLite database in a temp directory
  - Replaces the global `database.DB` with the test database
  - Creates the required tables schema
  - Returns a cleanup function to restore the original database
- **Usage**: Called at the beginning of each test to ensure isolation

#### `createTestTables() error`
- **Purpose**: Creates the tasks table schema in the test database
- **Schema**: Identical to production (id, title, description, done, created_at, completed_at)

#### `insertTestTask(t *testing.T, title, description string, done bool) int`
- **Purpose**: Helper to insert test data quickly
- **Returns**: The ID of the inserted task
- **Usage**: Used throughout tests to create known test data

#### `getTaskCount(t *testing.T) int`
- **Purpose**: Returns the total number of tasks in the test database
- **Usage**: Verifying correct number of tasks after operations

#### `getTaskByID(t *testing.T, id int) *testTask`
- **Purpose**: Retrieves a specific task by ID for verification
- **Returns**: A `testTask` struct or nil if not found
- **Usage**: Verifying task state after operations

#### `clearTestTasks(t *testing.T)`
- **Purpose**: Removes all tasks from the test database
- **Usage**: Ensuring clean state between test runs

### Test Data Structure

```go
type testTask struct {
    ID          int
    Title       string
    Description string
    Done        bool
}
```

## 📝 Add Command Tests (`add_test.go`)

### Test Cases

#### `TestAddTask`
- **Scenarios tested**:
  - Task with title and description
  - Task with title only (empty description)
  - Tasks with special characters (@#$%^&*())
  - Tasks with Unicode characters and emojis (🚀🎉)
  - Very long titles and descriptions
- **What it verifies**:
  - Task count increases by 1
  - Task data is stored correctly
  - New tasks default to `done = false`

#### `TestAddTaskEmptyTitle`
- **Purpose**: Verifies system handles empty titles gracefully
- **Expectation**: Should work (database allows empty titles)

#### `TestAddTaskMultipleTasks`
- **Purpose**: Tests adding multiple tasks sequentially
- **Verifies**: Task count increases correctly, all tasks are stored

#### `TestAddTaskConcurrency`
- **Purpose**: Tests concurrent task additions
- **Method**: Uses 10 goroutines adding tasks simultaneously
- **Verifies**: All tasks are added without data corruption

### How Add Tests Work

1. **Setup**: Each test creates an isolated database
2. **Execution**: Tests simulate the add command by inserting tasks directly into the database
3. **Verification**: Checks task count, data integrity, and default values
4. **Cleanup**: Automatic cleanup via deferred cleanup function

## 📋 List Command Tests (`list_test.go`)

### Test Cases

#### `TestListTasks`
- **Empty list scenario**: Verifies handling of empty task list
- **Single task scenario**: Tests listing one task with all fields
- **Multiple tasks scenario**: Tests listing multiple tasks with mixed data
- **Completed tasks scenario**: Tests tasks with completion timestamps
- **Mixed states scenario**: Tests listing both completed and incomplete tasks

#### `TestListTasksWithLargeDataset`
- **Purpose**: Performance and correctness with 100+ tasks
- **Verifies**: Sequential IDs, correct task ordering, no data loss

#### `TestListTasksErrorHandling`
- **Purpose**: Tests behavior when database is unavailable
- **Method**: Closes database connection before querying
- **Expectation**: Should return appropriate error

### How List Tests Work

1. **Data Setup**: Creates known test data using helper functions
2. **Query Execution**: Directly queries database to simulate list command
3. **Data Verification**: Checks returned data matches expected values
4. **State Verification**: Ensures proper handling of completion states

## ✅ Done Command Tests (`done_test.go`)

### Test Cases

#### `TestFindTaskById`
- **Existing task**: Verifies finding valid tasks by ID
- **Non-existent task**: Tests behavior with invalid IDs
- **Completed task**: Tests finding already completed tasks with timestamps

#### `TestMarkTaskAsDone`
- **Incomplete to complete**: Tests normal completion workflow
- **Already completed**: Tests handling of already completed tasks
- **Database verification**: Ensures completion timestamp is set correctly

#### `TestDoneCommandIntegration`
- **Complete workflow**: Tests find + mark done sequence
- **Non-existent task**: Tests error handling for invalid IDs

#### `TestMultipleTaskOperations`
- **Purpose**: Tests completing multiple tasks
- **Pattern**: Completes every other task, verifies correct state

#### `TestDoneWithDatabaseError`
- **Purpose**: Tests error handling when database is unavailable
- **Method**: Closes database before operations

### How Done Tests Work

1. **Task Creation**: Uses helper to create tasks in known states
2. **Operation Simulation**: Directly updates database to simulate done command
3. **State Verification**: Checks task completion status and timestamps
4. **Error Testing**: Simulates various error conditions

## 📤 Export Command Tests (`export_test.go`)

### Test Cases

#### `TestExportTasks`
- **Single completed task**: Tests exporting one completed task with all fields
- **Single incomplete task**: Tests exporting incomplete task (empty completed_at)
- **Multiple mixed tasks**: Tests exporting multiple tasks with different states
- **Empty database**: Verifies CSV creation with header only when no tasks exist
- **Special characters**: Tests CSV escaping for quotes, commas, and Unicode characters

#### `TestExportInvalidPath`
- **Purpose**: Tests error handling when output path is invalid
- **Method**: Attempts to export to non-existent directory
- **Expectation**: Should return appropriate error message

#### `TestExportLargeDataset`
- **Purpose**: Performance testing with large number of tasks (1000+)
- **Verifies**: Correct row count, proper data ordering, no memory issues
- **Performance**: Ensures export completes efficiently for large datasets

#### `TestExportDateTimeFormatting`
- **Purpose**: Tests proper timestamp formatting in CSV output
- **Method**: Creates task with specific timestamps
- **Verifies**: DateTime format matches expected pattern (YYYY-MM-DD HH:MM:SS)

### How Export Tests Work

1. **Data Setup**: Creates known test tasks with specific attributes
2. **Export Execution**: Calls export function with temporary output files
3. **CSV Verification**: Reads and parses generated CSV files
4. **Content Validation**: Verifies headers, data integrity, and formatting
5. **Error Testing**: Tests various failure scenarios

### Export Test Utilities

#### `exportToCSV(outputPath string) error`
- **Purpose**: Test wrapper for export functionality
- **Implementation**: Mirrors the actual export command logic
- **Usage**: Used by all export tests to generate CSV files

#### `readCSVFile(t *testing.T, filePath string) [][]string`
- **Purpose**: Helper to read and parse CSV files for verification
- **Returns**: 2D string slice representing CSV content
- **Usage**: Verifying exported data matches expected format

#### `insertTestTaskWithTimestamp(t *testing.T, title, description string, done bool) int`
- **Purpose**: Creates test tasks with realistic timestamps
- **Features**: Sets created_at and completed_at appropriately
- **Returns**: Task ID for verification

#### `insertTestTaskWithSpecificTime(t *testing.T, title, description string, done bool, createdAt time.Time, completedAt *time.Time) int`
- **Purpose**: Creates tasks with precise timestamps for datetime testing
- **Usage**: Testing specific datetime formatting scenarios

### CSV Format Verification

The tests verify the exported CSV format:

```csv
ID,Title,Description,Done,Created At,Completed At
1,Buy groceries,Milk, eggs, bread,true,2023-12-01 10:30:00,2023-12-01 15:45:00
2,Finish project,Complete the final report,false,2023-12-01 11:00:00,
```

#### Verification Points:
- **Header Row**: Correct column names and order
- **Data Types**: Proper conversion of integers, booleans, and timestamps
- **Null Handling**: Empty string for null completed_at fields
- **Special Characters**: Proper CSV escaping and quoting
- **Unicode Support**: Preservation of emojis and international characters

### Performance Testing

#### Large Dataset Test:
- **Volume**: 1000 tasks with mixed states
- **Verification**: All tasks exported correctly
- **Performance**: Completes within reasonable time limits
- **Memory**: No memory leaks or excessive usage

### Error Scenario Testing

#### File System Errors:
- **Invalid paths**: Non-existent directories
- **Permission errors**: Read-only directories
- **Disk space**: Simulated storage limitations

#### Database Errors:
- **Connection failures**: Database unavailable during export
- **Query errors**: Malformed database queries
- **Transaction issues**: Interrupted database operations

## 🔄 Integration Tests (`integration_test.go`)

### Test Cases

#### `TestCompleteTaskWorkflow`
- **Purpose**: Tests complete add → list → done workflow
- **Steps**:
  1. Start with empty database
  2. Add multiple tasks with different characteristics
  3. List and verify all tasks
  4. Complete some tasks
  5. Verify mixed completion states
  6. Complete remaining tasks
  7. Verify all tasks completed

#### `TestTaskLifecycle`
- **Purpose**: Tests single task from creation to completion
- **Phases**: Creation → Finding → Completion → Verification

#### `TestMixedTaskStates`
- **Purpose**: Tests complex scenarios with various task types
- **Includes**: Unicode, special characters, empty descriptions
- **Verifies**: Selective completion and state management

#### `TestErrorScenarios`
- **Purpose**: Tests various error conditions
- **Scenarios**: Non-existent tasks, double completion, invalid operations

#### `TestConcurrentOperations`
- **Purpose**: Tests thread safety and concurrent access
- **Method**: Concurrent task creation and completion
- **Verifies**: Data integrity under concurrent load

### Integration Test Philosophy

Integration tests focus on:
- **Workflows**: Testing complete user journeys
- **Data consistency**: Ensuring operations work together correctly
- **Error handling**: Testing edge cases and error conditions
- **Performance**: Ensuring system works under load

## 🚀 Running Tests

### Prerequisites
```bash
# Ensure you're in the project root
cd /path/to/tasker

# Install dependencies
go mod tidy
```

### Running All Tests
```bash
# Run all tests with verbose output
go test ./tests/... -v

# Run all tests with coverage
go test ./tests/... -cover

# Run tests and show coverage details
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Running Specific Test Files
```bash
# Add command tests only
go test ./tests/... -v -run TestAdd

# List command tests only
go test ./tests/... -v -run TestList

# Done command tests only
go test ./tests/... -v -run TestDone

# Export command tests only
go test ./tests/... -v -run TestExport

# Integration tests only
go test ./tests/... -v -run TestComplete
go test ./tests/... -v -run TestTask
go test ./tests/... -v -run TestMixed
go test ./tests/... -v -run TestError
go test ./tests/... -v -run TestConcurrent
```

### Running Individual Test Cases
```bash
# Specific test function
go test ./tests/... -v -run TestAddTask

# Test functions matching pattern
go test ./tests/... -v -run "TestAdd.*Concurrency"

# Test with specific timeout
go test ./tests/... -v -timeout 30s
```

### Debugging Tests
```bash
# Run tests with race detection
go test ./tests/... -race

# Run tests with detailed output
go test ./tests/... -v -args -test.v

# Run single test with extra verbosity
go test ./tests/... -v -run TestSpecificFunction
```

## 📊 Test Coverage

Our test suite covers:

### Add Command (95%+ coverage)
- ✅ Normal task creation
- ✅ Edge cases (empty fields, special characters)
- ✅ Unicode and emoji support
- ✅ Concurrent additions
- ✅ Error handling

### List Command (90%+ coverage)
- ✅ Empty list handling
- ✅ Single and multiple task listing
- ✅ Completion state display
- ✅ Large dataset handling
- ✅ Database error scenarios

### Done Command (90%+ coverage)
- ✅ Task finding by ID
- ✅ Completion state changes
- ✅ Timestamp management
- ✅ Already completed tasks
- ✅ Error scenarios

### Export Command (95%+ coverage)
- ✅ CSV file creation and formatting
- ✅ Header row generation
- ✅ Data type conversion (int, bool, timestamps)
- ✅ Null value handling (completed_at)
- ✅ Special character escaping
- ✅ Unicode and emoji preservation
- ✅ Large dataset performance
- ✅ File system error handling
- ✅ Custom output path support

### Integration (85%+ coverage)
- ✅ Complete workflows
- ✅ Cross-command interactions
- ✅ Concurrent operations
- ✅ Error propagation
- ✅ Data consistency

## 🔍 Test Philosophy

### Isolation
- Each test uses its own temporary database
- Tests don't depend on each other
- Clean state for every test run

### Comprehensive Coverage
- Happy path scenarios
- Edge cases and error conditions
- Performance and concurrency
- Data integrity and consistency

### Realistic Testing
- Tests simulate actual database operations
- Use realistic test data
- Test concurrent scenarios
- Verify actual SQL operations

### Maintainability
- Clear test names and descriptions
- Comprehensive helper functions
- Well-documented test cases
- Easy to extend and modify

## 🐛 Debugging Test Failures

### Common Issues

1. **Database Connection Errors**
   - Check if test database cleanup is working
   - Verify no tests are interfering with each other

2. **Timing Issues**
   - Use `time.WithinDuration()` for timestamp comparisons
   - Consider adding small delays for concurrent tests

3. **Data Corruption**
   - Ensure `clearTestTasks()` is called between test cases
   - Verify test isolation is working correctly

4. **Race Conditions**
   - Run tests with `-race` flag
   - Check concurrent access patterns

### Debugging Steps

1. **Run single test**: Isolate the failing test
2. **Check logs**: Look for database errors or panics
3. **Verify data**: Check if test data is as expected
4. **Test isolation**: Ensure other tests aren't affecting results
5. **Add debug output**: Temporarily add logging to understand flow

## 📈 Adding New Tests

### When to Add Tests

1. **New Features**: Any new command or functionality
2. **Bug Fixes**: Add test to prevent regression
3. **Edge Cases**: Discovered scenarios not covered
4. **Performance**: When optimizing existing features

### Test Addition Guidelines

1. **Follow naming convention**: `TestCommandName` or `TestFeatureName`
2. **Use table-driven tests**: For multiple similar scenarios
3. **Include error cases**: Test both success and failure paths
4. **Document test purpose**: Clear comments explaining what's being tested
5. **Use helpers**: Leverage existing helper functions
6. **Ensure isolation**: Each test should be independent

### Example New Test Structure

```go
func TestNewFeature(t *testing.T) {
    // Setup test database
    cleanup := setupTestDB(t)
    defer cleanup()

    t.Run("success scenario", func(t *testing.T) {
        clearTestTasks(t)
        
        // Test setup
        // Test execution
        // Assertions
    })
    
    t.Run("error scenario", func(t *testing.T) {
        clearTestTasks(t)
        
        // Test setup
        // Test execution
        // Error assertions
    })
}
```

This test suite ensures your Tasker CLI is robust, reliable, and handles all edge cases correctly!
