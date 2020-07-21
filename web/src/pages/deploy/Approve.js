import React from 'react';
import { Modal, Form, Input, Switch, message } from 'antd';

class Approve extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
    }
  }

  handleSubmit = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        this.setState({loading: true});
        values.IsPass    =  values.IsPass ? 1 : 0
        values.id        =  this.props.id
        dispatch({
          type: 'deploy/deployReview',
          payload: values,
        }).then(() => {
          this.setState({loading: false});
          this.props.approveCanael()
        })
      }
    });
  };

  render() {
    const {getFieldDecorator, getFieldValue} = this.props.form;
    return (
      <Modal
        visible
        width={600}
        maskClosable={false}
        title="审核发布申请"
        onCancel={this.props.approveCanael}
        confirmLoading={this.state.loading}
        onOk={this.handleSubmit}>
        <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
          <Form.Item label="审批结果">
            {getFieldDecorator('IsPass', {
              initialValue: true, 
              valuePropName: "checked",
              rules: [{ required: true }],
            })(
              <Switch checkedChildren="通过" unCheckedChildren="驳回"/>
            )}
          </Form.Item>
          <Form.Item label={getFieldValue('IsPass') ? '审批意见' : '驳回原因'}>
            {getFieldDecorator('Reason', {
              rules: [{ required: getFieldValue('IsPass') == false }],
            })(
              <Input.TextArea placeholder={getFieldValue('IsPass') ? '请输入审批意见' : '请输入驳回原因'}/>
            )}
          </Form.Item>
        </Form>
      </Modal>
    )
  }
}

export default Form.create()(Approve)
